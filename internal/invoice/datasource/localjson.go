package datasource

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"biz/internal/invoice"
	perr "biz/internal/platform/errors"
)

type LocalJSON struct{}

type payload struct {
	Invoices []invoice.InvoiceDraft `json:"invoices"`
}

func (LocalJSON) load(path string) (payload, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return payload{}, perr.Wrap(perr.KindDependencyUnavailable, "failed to read local source file", err)
	}
	var p payload
	if err := json.Unmarshal(b, &p); err == nil && len(p.Invoices) > 0 {
		return p, nil
	}
	var single invoice.InvoiceDraft
	if err := json.Unmarshal(b, &single); err != nil {
		return payload{}, perr.Wrap(perr.KindValidation, "invalid local json format", err)
	}
	p.Invoices = []invoice.InvoiceDraft{single}
	return p, nil
}

func (l LocalJSON) GetByID(_ context.Context, path, id string) (invoice.InvoiceDraft, error) {
	p, err := l.load(path)
	if err != nil {
		return invoice.InvoiceDraft{}, err
	}
	for _, inv := range p.Invoices {
		if inv.PageID == id || inv.InvoiceNumber == id {
			return inv, nil
		}
	}
	return invoice.InvoiceDraft{}, perr.New(perr.KindNotFound, "invoice not found in local source")
}

func (l LocalJSON) ListByStatus(_ context.Context, path, status string, limit int, _ string) ([]invoice.InvoiceSummary, string, error) {
	p, err := l.load(path)
	if err != nil {
		return nil, "", err
	}
	items := make([]invoice.InvoiceSummary, 0, len(p.Invoices))
	for _, inv := range p.Invoices {
		if status != "" && !strings.EqualFold(inv.Status, status) {
			continue
		}
		total := 0.0
		for _, li := range inv.LineItems {
			total += (li.Quantity * li.UnitRate) - li.Discount
		}
		items = append(items, invoice.InvoiceSummary{
			PageID:        inv.PageID,
			InvoiceNumber: inv.InvoiceNumber,
			ClientName:    inv.ClientName,
			Status:        inv.Status,
			DueDate:       inv.DueDate,
			Currency:      inv.Currency,
			Total:         total,
		})
	}
	if limit <= 0 || len(items) <= limit {
		return items, "", nil
	}
	return items[:limit], "local:next", nil
}
