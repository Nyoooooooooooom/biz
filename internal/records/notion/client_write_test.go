package notion

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientCreateUpdateArchive(t *testing.T) {
	t.Parallel()

	var (
		created  bool
		updated  bool
		archived bool
		schemaed bool
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer token-1" {
			t.Fatalf("unexpected auth header: %q", got)
		}
		if got := r.Header.Get("Notion-Version"); got != notionVersion {
			t.Fatalf("unexpected notion version: %q", got)
		}

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/databases/db_123":
			schemaed = true
			_, _ = io.WriteString(w, `{"id":"db_123","properties":{"Status":{"type":"status"},"Notes":{"type":"rich_text"}}}`)
			return
		case r.Method == http.MethodPost && r.URL.Path == "/v1/pages":
			created = true
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode create body: %v", err)
			}
			parent, _ := body["parent"].(map[string]any)
			if parent["database_id"] != "db_123" {
				t.Fatalf("unexpected database id: %#v", parent)
			}
			_, _ = io.WriteString(w, `{"id":"pg_1","last_edited_time":"2026-01-01T00:00:00Z","properties":{"Name":{"title":[{"plain_text":"A"}]}}}`)
			return
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/pages/pg_1":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode patch body: %v", err)
			}
			if archivedRaw, ok := body["archived"]; ok {
				archived = true
				if archivedRaw != true {
					t.Fatalf("expected archived=true, got %#v", archivedRaw)
				}
				_, _ = io.WriteString(w, `{"id":"pg_1","last_edited_time":"2026-01-01T00:00:00Z","properties":{}}`)
				return
			}
			updated = true
			_, _ = io.WriteString(w, `{"id":"pg_1","last_edited_time":"2026-01-01T00:00:00Z","properties":{"Status":{"status":{"name":"Done"}}}}`)
			return
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	c := New("token-1", srv.URL+"/v1", 2*time.Second, 1, 10)
	if _, err := c.Schema(context.Background(), "db_123"); err != nil {
		t.Fatalf("Schema failed: %v", err)
	}
	if _, err := c.Create(context.Background(), "db_123", map[string]any{"Name": map[string]any{"title": []any{map[string]any{"text": map[string]any{"content": "A"}}}}}); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if _, err := c.Update(context.Background(), "pg_1", map[string]any{"Status": map[string]any{"status": map[string]any{"name": "Done"}}}); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if err := c.Archive(context.Background(), "pg_1"); err != nil {
		t.Fatalf("Archive failed: %v", err)
	}
	if !schemaed || !created || !updated || !archived {
		t.Fatalf("expected schema/create/update/archive calls, got schema=%v created=%v updated=%v archived=%v", schemaed, created, updated, archived)
	}
}
