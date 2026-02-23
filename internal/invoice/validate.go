package invoice

import (
	"context"
	"math"
	"strings"

	perr "biz/internal/platform/errors"
)

func (s *Service) Validate(_ context.Context, in InvoiceDraft) error {
	if strings.TrimSpace(in.PageID) == "" {
		return perr.New(perr.KindValidation, "page_id is required")
	}
	if strings.TrimSpace(in.InvoiceNumber) == "" {
		return perr.New(perr.KindValidation, "invoice_number is required")
	}
	if in.IssueDate.IsZero() || in.DueDate.IsZero() {
		return perr.New(perr.KindValidation, "issue_date and due_date are required")
	}
	if in.DueDate.Before(in.IssueDate) {
		return perr.New(perr.KindValidation, "due_date must be on/after issue_date")
	}
	if strings.TrimSpace(in.Currency) == "" {
		return perr.New(perr.KindValidation, "currency is required")
	}
	allowed := false
	for _, c := range s.CurrencyAllowList {
		if strings.EqualFold(c, in.Currency) {
			allowed = true
			break
		}
	}
	if !allowed {
		return perr.New(perr.KindValidation, "currency is not allowed")
	}
	if len(in.LineItems) == 0 {
		return perr.New(perr.KindValidation, "at least one line item is required")
	}
	for i, li := range in.LineItems {
		if strings.TrimSpace(li.Description) == "" {
			return perr.New(perr.KindValidation, "line item description is required")
		}
		if li.Quantity <= 0 {
			return perr.New(perr.KindValidation, "line item quantity must be > 0 at index "+itoa(i))
		}
		if li.UnitRate < 0 || li.Discount < 0 {
			return perr.New(perr.KindValidation, "line item rate/discount must be >= 0 at index "+itoa(i))
		}
		if li.Discount > li.UnitRate*li.Quantity {
			return perr.New(perr.KindValidation, "line item discount exceeds line amount at index "+itoa(i))
		}
	}
	if t := s.computeSubtotal(in); t < 0 || math.IsNaN(t) {
		return perr.New(perr.KindValidation, "invalid totals")
	}
	return nil
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	buf := make([]byte, 0, 8)
	for v > 0 {
		buf = append([]byte{byte('0' + (v % 10))}, buf...)
		v /= 10
	}
	return string(buf)
}
