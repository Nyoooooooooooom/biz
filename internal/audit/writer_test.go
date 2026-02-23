package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"biz/internal/platform/config"
)

func TestWriterHashChain(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "audit.log")
	w, err := NewWriter(config.AuditConfig{Path: path, SigningKey: "test-key", DirPerm: 0o700, FilePerm: 0o600})
	if err != nil {
		t.Fatalf("new writer: %v", err)
	}
	if err := w.Write(Event{TraceID: "t1", Actor: "agent", Command: "invoice.list", ExitCode: 0, ResultCode: "OK"}); err != nil {
		t.Fatalf("write first: %v", err)
	}
	if err := w.Write(Event{TraceID: "t2", Actor: "agent", Command: "invoice.preview", ExitCode: 2, ResultCode: "VALIDATION_ERROR"}); err != nil {
		t.Fatalf("write second: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open log: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var rows []map[string]any
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		var row map[string]any
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			t.Fatalf("unmarshal row: %v", err)
		}
		rows = append(rows, row)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[1]["prev_hash"] != rows[0]["hash"] {
		t.Fatalf("expected hash chain linkage")
	}
	if rows[0]["signature"] == "" || rows[1]["signature"] == "" {
		t.Fatalf("expected non-empty signatures")
	}
}
