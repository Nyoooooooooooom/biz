package notion

import (
	"testing"
	"time"
)

func TestMapInvoicePage(t *testing.T) {
	invoicePage := rawPage{
		ID:             "page-1",
		LastEditedTime: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
		Properties: map[string]any{
			"Invoice Number":  map[string]any{"title": []any{map[string]any{"plain_text": "INV-001"}}},
			"Client Name":     map[string]any{"rich_text": []any{map[string]any{"plain_text": "Acme Corp"}}},
			"Client Location": map[string]any{"select": map[string]any{"name": "US"}},
			"Invoice Date":    map[string]any{"date": map[string]any{"start": "2026-01-01"}},
			"Due Date":        map[string]any{"date": map[string]any{"start": "2026-01-10"}},
			"Currency":        map[string]any{"select": map[string]any{"name": "usd"}},
			"Status":          map[string]any{"select": map[string]any{"name": "Ready to Invoice"}},
		},
	}
	worklogPage := rawPage{Properties: map[string]any{
		"Description":    map[string]any{"rich_text": []any{map[string]any{"plain_text": "Consulting"}}},
		"Hours":          map[string]any{"number": 2.0},
		"Effective Rate": map[string]any{"formula": map[string]any{"number": 150.0}},
	}}

	d, err := mapInvoicePage(invoicePage, []rawPage{worklogPage}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.InvoiceNumber != "INV-001" || d.Currency != "USD" {
		t.Fatalf("unexpected mapped invoice: %+v", d)
	}
	if len(d.LineItems) != 1 || d.LineItems[0].UnitRate != 150 || d.LineItems[0].Quantity != 2 {
		t.Fatalf("unexpected line items: %+v", d.LineItems)
	}
}

func TestMapInvoicePageWithWorklogsAndCosts(t *testing.T) {
	invoicePage := rawPage{
		ID:             "inv-1",
		LastEditedTime: time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC),
		Properties: map[string]any{
			"Invoice Number": map[string]any{"title": []any{map[string]any{"plain_text": "INV-100"}}},
			"Invoice Date":   map[string]any{"date": map[string]any{"start": "2026-02-01"}},
			"Due Date":       map[string]any{"date": map[string]any{"start": "2026-02-15"}},
			"Currency":       map[string]any{"select": map[string]any{"name": "usd"}},
		},
	}
	worklogPage := rawPage{
		Properties: map[string]any{
			"Description":    map[string]any{"rich_text": []any{map[string]any{"plain_text": "Flowcal fix"}}},
			"Hours":          map[string]any{"number": 1.5},
			"Effective Rate": map[string]any{"formula": map[string]any{"number": 200.0}},
		},
	}
	costPage := rawPage{
		Properties: map[string]any{
			"Name":            map[string]any{"title": []any{map[string]any{"plain_text": "Hosting"}}},
			"Category":        map[string]any{"select": map[string]any{"name": "Hosting"}},
			"Billable Amount": map[string]any{"formula": map[string]any{"number": 50.0}},
		},
	}

	clientPage := &rawPage{
		Properties: map[string]any{
			"Name":       map[string]any{"title": []any{map[string]any{"plain_text": "Acme Industries"}}},
			"Tax Region": map[string]any{"select": map[string]any{"name": "US"}},
		},
	}
	d, err := mapInvoicePage(invoicePage, []rawPage{worklogPage}, []rawPage{costPage}, clientPage)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d.LineItems) != 2 {
		t.Fatalf("expected 2 line items, got %d", len(d.LineItems))
	}
	if d.ClientName != "Acme Industries" || d.ClientLocation != "US" {
		t.Fatalf("unexpected client mapping: %+v", d)
	}
	if d.LineItems[0].Quantity != 1.5 || d.LineItems[0].UnitRate != 200 {
		t.Fatalf("unexpected worklog item: %+v", d.LineItems[0])
	}
	if d.LineItems[1].Quantity != 1 || d.LineItems[1].UnitRate != 50 {
		t.Fatalf("unexpected cost item: %+v", d.LineItems[1])
	}
}

func TestRelationIDs(t *testing.T) {
	props := map[string]any{
		"Line Items": map[string]any{
			"relation": []any{map[string]any{"id": "a"}, map[string]any{"id": "b"}},
		},
	}
	ids := relationIDs(props, "Line Items")
	if len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Fatalf("unexpected ids: %+v", ids)
	}
}
