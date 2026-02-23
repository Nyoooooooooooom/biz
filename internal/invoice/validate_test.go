package invoice

import (
	"context"
	"testing"
	"time"
)

func TestValidateInvoiceDraft(t *testing.T) {
	s := &Service{CurrencyAllowList: []string{"USD", "EUR"}}
	d := InvoiceDraft{
		PageID:        "p1",
		InvoiceNumber: "INV-1",
		ClientName:    "Acme",
		IssueDate:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		DueDate:       time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
		Currency:      "USD",
		LineItems: []LineItem{{
			Description: "Consulting",
			Quantity:    1,
			UnitRate:    100,
		}},
	}
	if err := s.Validate(context.Background(), d); err != nil {
		t.Fatalf("expected valid draft, got: %v", err)
	}
	d.Currency = "JPY"
	if err := s.Validate(context.Background(), d); err == nil {
		t.Fatalf("expected currency validation error")
	}
}
