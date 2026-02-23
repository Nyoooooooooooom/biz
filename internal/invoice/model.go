package invoice

import "time"

type LineItem struct {
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitRate    float64 `json:"unit_rate"`
	Discount    float64 `json:"discount,omitempty"`
	SortOrder   int     `json:"sort_order,omitempty"`
}

type InvoiceDraft struct {
	PageID         string     `json:"page_id"`
	InvoiceNumber  string     `json:"invoice_number"`
	ClientName     string     `json:"client_name"`
	ClientLocation string     `json:"client_location"`
	IssueDate      time.Time  `json:"issue_date"`
	DueDate        time.Time  `json:"due_date"`
	Currency       string     `json:"currency"`
	Status         string     `json:"status"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	LineItems      []LineItem `json:"line_items"`
	Notes          string     `json:"notes,omitempty"`
}

type Totals struct {
	Subtotal float64 `json:"subtotal"`
	TaxRate  float64 `json:"tax_rate"`
	Tax      float64 `json:"tax"`
	Total    float64 `json:"total"`
}

type InvoiceDocument struct {
	Draft       InvoiceDraft `json:"draft"`
	Totals      Totals       `json:"totals"`
	Idempotency string       `json:"idempotency_key"`
}

type CreateRequest struct {
	ID           string
	OutDir       string
	UploadNotion bool
	Confirm      bool
	Source       string
	SourceFile   string
	TraceID      string
}

type CreateResult struct {
	PageID         string         `json:"page_id"`
	InvoiceNumber  string         `json:"invoice_number"`
	PDFPath        string         `json:"pdf_path"`
	PDFSizeBytes   int64          `json:"pdf_size_bytes"`
	IdempotencyKey string         `json:"idempotency_key"`
	Totals         Totals         `json:"totals"`
	Meta           map[string]any `json:"meta,omitempty"`
}

type ListRequest struct {
	Status  string
	Limit   int
	Cursor  string
	TraceID string
}

type InvoiceSummary struct {
	PageID        string    `json:"page_id"`
	InvoiceNumber string    `json:"invoice_number"`
	ClientName    string    `json:"client_name"`
	Status        string    `json:"status"`
	DueDate       time.Time `json:"due_date"`
	Total         float64   `json:"total"`
	Currency      string    `json:"currency"`
}

type ListResult struct {
	Items      []InvoiceSummary `json:"items"`
	NextCursor string           `json:"next_cursor,omitempty"`
}

type PreviewRequest struct {
	ID      string
	Format  string
	TraceID string
}

type PreviewResult struct {
	Path      string `json:"path"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}
