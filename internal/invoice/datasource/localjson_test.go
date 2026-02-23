package datasource

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalJSONGetByIDAndList(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "invoices.json")
	payload := `{
  "invoices": [
    {
      "page_id": "p1",
      "invoice_number": "INV-1",
      "client_name": "Acme",
      "client_location": "US",
      "issue_date": "2026-01-01T00:00:00Z",
      "due_date": "2026-01-10T00:00:00Z",
      "currency": "USD",
      "status": "Overdue",
      "last_edited_time": "2026-01-05T00:00:00Z",
      "line_items": [{"description":"Work","quantity":1,"unit_rate":100}]
    }
  ]
}`
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	src := LocalJSON{}
	got, err := src.GetByID(context.Background(), path, "p1")
	if err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if got.InvoiceNumber != "INV-1" {
		t.Fatalf("unexpected invoice: %+v", got)
	}
	items, _, err := src.ListByStatus(context.Background(), path, "Overdue", 10, "")
	if err != nil {
		t.Fatalf("ListByStatus error: %v", err)
	}
	if len(items) != 1 || items[0].Total != 100 {
		t.Fatalf("unexpected list: %+v", items)
	}
}
