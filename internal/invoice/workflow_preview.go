package invoice

import (
	"context"
	"os"
	"path/filepath"

	perr "biz/internal/platform/errors"
	"biz/internal/tax"
)

func (s *Service) workflowPreview(ctx context.Context, req PreviewRequest) (PreviewResult, error) {
	if req.Format != "" && req.Format != "html" && req.Format != "pdf" {
		return PreviewResult{}, perr.New(perr.KindValidation, "format must be html or pdf")
	}
	draft, _, err := s.fetchDraft(ctx, req.ID, s.resolveSource(""), "")
	if err != nil {
		return PreviewResult{}, err
	}
	if err := s.Validate(ctx, draft); err != nil {
		return PreviewResult{}, err
	}
	sub := s.computeSubtotal(draft)
	taxOut, err := s.Tax.Apply(ctx, tax.TaxInput{Region: draft.ClientLocation, Subtotal: sub, Currency: draft.Currency})
	if err != nil {
		return PreviewResult{}, err
	}
	doc := InvoiceDocument{
		Draft:       draft,
		Totals:      Totals{Subtotal: sub, TaxRate: taxOut.Rate, Tax: taxOut.Amount, Total: round2(sub + taxOut.Amount)},
		Idempotency: "preview",
	}
	htmlBytes, err := s.Template.RenderInvoiceHTML(ctx, doc)
	if err != nil {
		return PreviewResult{}, perr.Wrap(perr.KindInternal, "failed to render preview html", err)
	}
	if s.MaxRenderHTMLBytes > 0 && len(htmlBytes) > s.MaxRenderHTMLBytes {
		return PreviewResult{}, perr.New(perr.KindValidation, "rendered invoice HTML exceeds configured safety limit")
	}

	if req.Format == "html" {
		p := filepath.Join(os.TempDir(), "biz-preview-"+req.ID+".html")
		if err := os.WriteFile(p, htmlBytes, 0o644); err != nil {
			return PreviewResult{}, perr.Wrap(perr.KindInternal, "failed to write preview html", err)
		}
		st, _ := os.Stat(p)
		return PreviewResult{Path: p, MimeType: "text/html", SizeBytes: st.Size()}, nil
	}
	pdfCtx, cancel := context.WithTimeout(ctx, s.PDFTimeout)
	defer cancel()
	pdfBytes, err := s.PDF.RenderPDF(pdfCtx, htmlBytes)
	if err != nil {
		return PreviewResult{}, perr.Wrap(perr.KindDependencyUnavailable, "failed to render preview pdf", err)
	}
	name := "preview_" + draft.InvoiceNumber + "_" + s.Clock.Now().Format("20060102T150405")
	p, size, err := s.LocalPDF.Save(ctx, name, draft.IssueDate, "preview", pdfBytes, os.TempDir())
	if err != nil {
		return PreviewResult{}, err
	}
	return PreviewResult{Path: p, MimeType: "application/pdf", SizeBytes: size}, nil
}
