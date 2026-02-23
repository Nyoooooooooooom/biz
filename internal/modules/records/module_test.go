package records

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePropertiesFromData(t *testing.T) {
	props, err := parseProperties(`{"Status":{"status":{"name":"Done"}}}`, "")
	if err != nil {
		t.Fatalf("parseProperties failed: %v", err)
	}
	if _, ok := props["Status"]; !ok {
		t.Fatalf("expected Status property, got %#v", props)
	}
}

func TestParsePropertiesFromFile(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "payload.json")
	if err := os.WriteFile(p, []byte(`{"Notes":{"rich_text":[{"text":{"content":"hi"}}]}}`), 0o600); err != nil {
		t.Fatalf("write payload: %v", err)
	}
	props, err := parseProperties("", p)
	if err != nil {
		t.Fatalf("parseProperties from file failed: %v", err)
	}
	if _, ok := props["Notes"]; !ok {
		t.Fatalf("expected Notes property, got %#v", props)
	}
}

func TestParsePropertiesRejectsBothDataAndFile(t *testing.T) {
	if _, err := parseProperties(`{"A":1}`, "/tmp/payload.json"); err == nil {
		t.Fatal("expected error when both --data and --data-file provided")
	}
}

func TestDiffProperties(t *testing.T) {
	current := map[string]any{
		"Status": map[string]any{"status": map[string]any{"name": "Todo"}},
		"Notes":  map[string]any{"rich_text": []any{}},
	}
	next := map[string]any{
		"Status": map[string]any{"status": map[string]any{"name": "Done"}},
		"Notes":  map[string]any{"rich_text": []any{}},
	}
	diff := diffProperties(current, next)
	if len(diff) != 1 {
		t.Fatalf("expected 1 changed property, got %d (%#v)", len(diff), diff)
	}
	if _, ok := diff["Status"]; !ok {
		t.Fatalf("expected Status diff, got %#v", diff)
	}
}
