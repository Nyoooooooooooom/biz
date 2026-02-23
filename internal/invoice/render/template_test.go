package render

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"biz/internal/invoice"
)

func TestRenderInvoiceHTML(t *testing.T) {
	tmp := t.TempDir()
	tp := filepath.Join(tmp, "invoice.tmpl")
	content := `<html><body>{{ .Draft.InvoiceNumber }} {{ printf "%.2f" (lineTotal 2 10 1) }}</body></html>`
	if err := os.WriteFile(tp, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	r := NewTemplate(tp)
	b, err := r.RenderInvoiceHTML(context.Background(), invoice.InvoiceDocument{
		Draft: invoice.InvoiceDraft{InvoiceNumber: "INV-1", IssueDate: time.Now(), DueDate: time.Now(), Currency: "USD", PageID: "p1", LineItems: []invoice.LineItem{{Description: "x", Quantity: 1, UnitRate: 1}}},
	})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if !strings.Contains(string(b), "INV-1") || !strings.Contains(string(b), "19.00") {
		t.Fatalf("unexpected output: %s", string(b))
	}
}
