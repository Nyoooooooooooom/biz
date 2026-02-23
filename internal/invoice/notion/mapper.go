package notion

import (
	"fmt"
	"sort"
	"strings"

	"biz/internal/invoice"
	perr "biz/internal/platform/errors"
)

func mapInvoicePage(p rawPage, worklogPages []rawPage, costPages []rawPage, clientPage *rawPage) (invoice.InvoiceDraft, error) {
	inv := invoice.InvoiceDraft{
		PageID:         p.ID,
		InvoiceNumber:  firstText(p.Properties, "Invoice Number"),
		ClientName:     firstText(p.Properties, "Client Name"),
		ClientLocation: firstText(p.Properties, "Client Location"),
		IssueDate:      firstDate(p.Properties, "Invoice Date"),
		DueDate:        firstDate(p.Properties, "Due Date"),
		Currency:       strings.ToUpper(firstText(p.Properties, "Currency")),
		Status:         firstText(p.Properties, "Status"),
		LastEditedTime: p.LastEditedTime,
		Notes:          firstText(p.Properties, "Notes"),
		LineItems:      make([]invoice.LineItem, 0, len(worklogPages)+len(costPages)),
	}
	if inv.PageID == "" {
		return invoice.InvoiceDraft{}, perr.New(perr.KindValidation, "invalid notion invoice payload")
	}
	if clientPage != nil {
		if inv.ClientName == "" {
			inv.ClientName = firstText(clientPage.Properties, "Name")
		}
		if inv.ClientLocation == "" {
			inv.ClientLocation = firstText(clientPage.Properties, "Tax Region")
		}
	}
	if inv.InvoiceNumber == "" || inv.IssueDate.IsZero() || inv.DueDate.IsZero() || inv.Currency == "" {
		return invoice.InvoiceDraft{}, perr.New(perr.KindValidation, "invoice page missing required properties: Invoice Number, Invoice Date, Due Date, Currency")
	}
	for i, wp := range worklogPages {
		inv.LineItems = append(inv.LineItems, mapWorklogPage(wp, i+1))
	}
	for i, cp := range costPages {
		inv.LineItems = append(inv.LineItems, mapCostPage(cp, len(inv.LineItems)+i+1))
	}
	sort.SliceStable(inv.LineItems, func(i, j int) bool {
		return inv.LineItems[i].SortOrder < inv.LineItems[j].SortOrder
	})
	return inv, nil
}

func mapSummaryPage(p rawPage) invoice.InvoiceSummary {
	return invoice.InvoiceSummary{
		PageID:        p.ID,
		InvoiceNumber: firstText(p.Properties, "Invoice Number"),
		ClientName:    firstText(p.Properties, "Client Name"),
		Status:        firstText(p.Properties, "Status"),
		DueDate:       firstDate(p.Properties, "Due Date"),
		Total:         firstNumber(p.Properties, "Total"),
		Currency:      strings.ToUpper(firstText(p.Properties, "Currency")),
	}
}

func mapWorklogPage(p rawPage, idx int) invoice.LineItem {
	desc := firstText(p.Properties, "Description", "Title")
	if desc == "" {
		desc = "Worklog"
	}
	hours := firstNumber(p.Properties, "Hours")
	if hours <= 0 {
		hours = firstNumber(p.Properties, "Minutes") / 60.0
	}
	rate := firstNumber(p.Properties, "Effective Rate", "Actual Rate", "Rate Override")
	if rate <= 0 && hours > 0 {
		amount := firstNumber(p.Properties, "Amount", "Total")
		if amount > 0 {
			rate = amount / hours
		}
	}
	return invoice.LineItem{
		Description: fmt.Sprintf("Worklog: %s", desc),
		Quantity:    hours,
		UnitRate:    rate,
		Discount:    0,
		SortOrder:   idx,
	}
}

func mapCostPage(p rawPage, idx int) invoice.LineItem {
	desc := firstText(p.Properties, "Name")
	if desc == "" {
		desc = "Cost"
	}
	category := firstText(p.Properties, "Category")
	if category != "" {
		desc = fmt.Sprintf("%s (%s)", desc, category)
	}
	amount := firstNumber(p.Properties, "Billable Amount")
	return invoice.LineItem{
		Description: fmt.Sprintf("Cost: %s", desc),
		Quantity:    1,
		UnitRate:    amount,
		Discount:    0,
		SortOrder:   idx,
	}
}
