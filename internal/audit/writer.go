package audit

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"biz/internal/platform/config"
	perr "biz/internal/platform/errors"
)

type Writer struct {
	path     string
	signing  []byte
	dirPerm  os.FileMode
	filePerm os.FileMode
	mu       sync.Mutex
}

type Event struct {
	TraceID      string
	Actor        string
	Command      string
	Args         []string
	ExitCode     int
	ResultCode   string
	ErrorMessage string
}

type record struct {
	Timestamp    string   `json:"timestamp"`
	TraceID      string   `json:"trace_id"`
	Actor        string   `json:"actor"`
	Command      string   `json:"command"`
	Args         []string `json:"args,omitempty"`
	ExitCode     int      `json:"exit_code"`
	ResultCode   string   `json:"result_code"`
	ErrorMessage string   `json:"error_message,omitempty"`
	PrevHash     string   `json:"prev_hash,omitempty"`
	Hash         string   `json:"hash"`
	Signature    string   `json:"signature"`
}

type unsignedRecord struct {
	Timestamp    string   `json:"timestamp"`
	TraceID      string   `json:"trace_id"`
	Actor        string   `json:"actor"`
	Command      string   `json:"command"`
	Args         []string `json:"args,omitempty"`
	ExitCode     int      `json:"exit_code"`
	ResultCode   string   `json:"result_code"`
	ErrorMessage string   `json:"error_message,omitempty"`
	PrevHash     string   `json:"prev_hash,omitempty"`
}

func NewWriter(cfg config.AuditConfig) (*Writer, error) {
	if strings.TrimSpace(cfg.Path) == "" {
		return nil, perr.New(perr.KindValidation, "audit.path is required when audit is enabled")
	}
	if strings.TrimSpace(cfg.SigningKey) == "" {
		return nil, perr.New(perr.KindValidation, "audit.signing_key is required when audit is enabled")
	}
	dirPerm := os.FileMode(cfg.DirPerm)
	if dirPerm == 0 {
		dirPerm = 0o700
	}
	filePerm := os.FileMode(cfg.FilePerm)
	if filePerm == 0 {
		filePerm = 0o600
	}
	return &Writer{
		path:     cfg.Path,
		signing:  []byte(cfg.SigningKey),
		dirPerm:  dirPerm,
		filePerm: filePerm,
	}, nil
}

func (w *Writer) Write(e Event) error {
	if w == nil {
		return nil
	}
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(w.path), w.dirPerm); err != nil {
		return perr.Wrap(perr.KindInternal, "failed to create audit dir", err)
	}
	prev, err := lastHash(w.path)
	if err != nil {
		return err
	}
	unsigned := unsignedRecord{
		Timestamp:    time.Now().UTC().Format(time.RFC3339Nano),
		TraceID:      strings.TrimSpace(e.TraceID),
		Actor:        strings.TrimSpace(e.Actor),
		Command:      strings.TrimSpace(e.Command),
		Args:         e.Args,
		ExitCode:     e.ExitCode,
		ResultCode:   strings.TrimSpace(e.ResultCode),
		ErrorMessage: strings.TrimSpace(e.ErrorMessage),
		PrevHash:     prev,
	}
	payload, err := json.Marshal(unsigned)
	if err != nil {
		return perr.Wrap(perr.KindInternal, "failed to marshal audit record", err)
	}
	hash := hashRecord(prev, payload)
	sig := sign(hash, w.signing)
	rec := record{
		Timestamp:    unsigned.Timestamp,
		TraceID:      unsigned.TraceID,
		Actor:        unsigned.Actor,
		Command:      unsigned.Command,
		Args:         unsigned.Args,
		ExitCode:     unsigned.ExitCode,
		ResultCode:   unsigned.ResultCode,
		ErrorMessage: unsigned.ErrorMessage,
		PrevHash:     unsigned.PrevHash,
		Hash:         hash,
		Signature:    sig,
	}
	line, err := json.Marshal(rec)
	if err != nil {
		return perr.Wrap(perr.KindInternal, "failed to marshal signed audit record", err)
	}
	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, w.filePerm)
	if err != nil {
		return perr.Wrap(perr.KindInternal, "failed to open audit log", err)
	}
	defer f.Close()
	if _, err := f.Write(append(line, '\n')); err != nil {
		return perr.Wrap(perr.KindInternal, "failed to write audit log", err)
	}
	if err := f.Chmod(w.filePerm); err != nil {
		return perr.Wrap(perr.KindInternal, "failed to enforce audit log file permissions", err)
	}
	return nil
}

func lastHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", perr.Wrap(perr.KindInternal, "failed to open audit log", err)
	}
	defer f.Close()

	var last string
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		var r struct {
			Hash string `json:"hash"`
		}
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			return "", perr.Wrap(perr.KindInternal, "failed to parse audit log line", err)
		}
		if strings.TrimSpace(r.Hash) == "" {
			return "", perr.New(perr.KindInternal, "invalid audit log: missing hash")
		}
		last = r.Hash
	}
	if err := s.Err(); err != nil {
		return "", perr.Wrap(perr.KindInternal, "failed to read audit log", err)
	}
	return last, nil
}

func hashRecord(prev string, payload []byte) string {
	sum := sha256.Sum256([]byte(prev + "|" + string(payload)))
	return hex.EncodeToString(sum[:])
}

func sign(hash string, key []byte) string {
	m := hmac.New(sha256.New, key)
	_, _ = m.Write([]byte(hash))
	return hex.EncodeToString(m.Sum(nil))
}

func (w *Writer) String() string {
	if w == nil {
		return ""
	}
	return fmt.Sprintf("audit(path=%s)", w.path)
}
