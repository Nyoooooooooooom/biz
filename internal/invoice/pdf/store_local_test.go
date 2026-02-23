package pdf

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveRejectsOutputOutsideBase(t *testing.T) {
	base := t.TempDir()
	store := LocalStore{OutputBaseDir: filepath.Join(base, "allowed")}
	_, _, err := store.Save(context.Background(), "INV-1", time.Now(), "abc123", []byte("pdf"), filepath.Join(base, "outside"))
	if err == nil {
		t.Fatal("expected error for output outside base dir")
	}
}
