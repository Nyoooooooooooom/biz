package pdf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	perr "biz/internal/platform/errors"
)

type LocalStore struct {
	OutputBaseDir string
	DirPerm       os.FileMode
	FilePerm      os.FileMode
}

func (s LocalStore) Save(_ context.Context, invoiceNumber string, issueDate time.Time, idempotencyKey string, pdf []byte, outDir string) (string, int64, error) {
	return s.save(invoiceNumber, issueDate, idempotencyKey, pdf, outDir)
}

func (s LocalStore) save(invoiceNumber string, issueDate time.Time, idempotencyKey string, pdf []byte, outDir string) (string, int64, error) {
	if outDir == "" {
		outDir = "invoices"
	}
	if s.DirPerm == 0 {
		s.DirPerm = 0o700
	}
	if s.FilePerm == 0 {
		s.FilePerm = 0o600
	}
	if err := ensureWithinBase(outDir, s.OutputBaseDir); err != nil {
		return "", 0, err
	}
	if err := os.MkdirAll(outDir, s.DirPerm); err != nil {
		return "", 0, perr.Wrap(perr.KindInternal, "failed to create output dir", err)
	}
	short := idempotencyKey
	if len(short) > 8 {
		short = short[:8]
	}
	name := fmt.Sprintf("%s_%s_%s.pdf", sanitize(invoiceNumber), issueDate.Format("2006-01-02"), short)
	path := filepath.Join(outDir, name)
	if err := os.WriteFile(path, pdf, s.FilePerm); err != nil {
		return "", 0, perr.Wrap(perr.KindInternal, "failed to write pdf", err)
	}
	st, err := os.Stat(path)
	if err != nil {
		return path, int64(len(pdf)), nil
	}
	return path, st.Size(), nil
}

func ensureWithinBase(outDir, baseDir string) error {
	if strings.TrimSpace(baseDir) == "" {
		return nil
	}
	outAbs, err := filepath.Abs(outDir)
	if err != nil {
		return perr.Wrap(perr.KindValidation, "invalid output directory", err)
	}
	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return perr.Wrap(perr.KindValidation, "invalid output base directory", err)
	}
	rel, err := filepath.Rel(baseAbs, outAbs)
	if err != nil {
		return perr.Wrap(perr.KindValidation, "invalid output directory relationship", err)
	}
	if rel == "." {
		return nil
	}
	if strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return perr.New(perr.KindValidation, "output directory must be inside configured output_base_dir")
	}
	return nil
}

func sanitize(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, " ", "-")
	v = strings.ReplaceAll(v, "/", "-")
	if v == "" {
		return "invoice"
	}
	return v
}
