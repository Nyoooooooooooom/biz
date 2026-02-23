# Notion Adapter

Path: `internal/invoice/notion`

## Purpose
Translate Notion pages/relations into `invoice.InvoiceDraft` and `invoice.InvoiceSummary`.

## Expected Invoice Schema
Invoice page properties:
- `Invoice Number` (title)
- `Invoice Date` (date)
- `Due Date` (date)
- `Currency` (select)
- `Status` (status)
- `Worklogs` (relation)
- `Costs` (relation)
- `Clients` (relation, optional for name/region enrichment)

Worklog page properties:
- `Description` (rich_text/title)
- `Hours` (number)
- `Minutes` (number, optional)
- `Effective Rate` or `Actual Rate` or `Rate Override`
- `Amount` or `Total` (fallback amount source)

Cost page properties:
- `Name` (title)
- `Category` (select)
- `Billable Amount` (number/formula)

## Notes
- `invoice list` requires `notion.invoice_db_id`.
- `invoice create/preview` can run from invoice page id + token.
