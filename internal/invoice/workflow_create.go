package invoice

import (
	"context"
	"fmt"
	"strings"

	perr "biz/internal/platform/errors"
	"biz/internal/tax"
)

func (s *Service) workflowCreate(ctx context.Context, req CreateRequest) (CreateResult, error) {
	if req.UploadNotion && !s.AllowMutations {
		return CreateResult{}, perr.New(perr.KindValidation, "notion mutations are disabled by policy")
	}
	if req.UploadNotion && s.RequireMutationConfirm && !req.Confirm {
		return CreateResult{}, perr.New(perr.KindValidation, "--confirm is required when --upload-notion is used")
	}

	source := s.resolveSource(req.Source)
	draft, usedFallback, err := s.fetchDraft(ctx, req.ID, source, req.SourceFile)
	if err != nil {
		return CreateResult{}, err
	}
	if err := s.Validate(ctx, draft); err != nil {
		return CreateResult{}, err
	}

	key := s.idempotencyKey(draft)
	store, err := s.loadIdempotency()
	if err != nil {
		return CreateResult{}, perr.Wrap(perr.KindInternal, "failed to load idempotency store", err)
	}
	if hit, ok := store[key]; ok {
		res := hit.Result
		if res.Meta == nil {
			res.Meta = map[string]any{}
		}
		res.Meta["idempotent_replay"] = true
		return res, nil
	}

	subtotal := s.computeSubtotal(draft)
	taxOut, err := s.Tax.Apply(ctx, tax.TaxInput{Region: draft.ClientLocation, Subtotal: subtotal, Currency: draft.Currency})
	if err != nil {
		return CreateResult{}, err
	}
	totals := Totals{Subtotal: subtotal, TaxRate: taxOut.Rate, Tax: taxOut.Amount, Total: round2(subtotal + taxOut.Amount)}
	if totals.Total < 0 {
		return CreateResult{}, perr.New(perr.KindValidation, "invoice total cannot be negative")
	}

	doc := InvoiceDocument{Draft: draft, Totals: totals, Idempotency: key}
	htmlBytes, err := s.Template.RenderInvoiceHTML(ctx, doc)
	if err != nil {
		return CreateResult{}, perr.Wrap(perr.KindInternal, "failed to render html", err)
	}
	if s.MaxRenderHTMLBytes > 0 && len(htmlBytes) > s.MaxRenderHTMLBytes {
		return CreateResult{}, perr.New(perr.KindValidation, "rendered invoice HTML exceeds configured safety limit")
	}

	pdfCtx, cancel := context.WithTimeout(ctx, s.PDFTimeout)
	defer cancel()
	pdfBytes, err := s.PDF.RenderPDF(pdfCtx, htmlBytes)
	if err != nil {
		return CreateResult{}, perr.Wrap(perr.KindDependencyUnavailable, "failed to render pdf", err)
	}
	outDir := req.OutDir
	if outDir == "" {
		outDir = "invoices"
	}
	path, size, err := s.LocalPDF.Save(ctx, draft.InvoiceNumber, draft.IssueDate, key, pdfBytes, outDir)
	if err != nil {
		if perr.KindOf(err) == perr.KindValidation {
			return CreateResult{}, err
		}
		return CreateResult{}, perr.Wrap(perr.KindInternal, "failed to store local pdf", err)
	}
	if req.UploadNotion {
		if err := s.NotionPDF.Store(ctx, draft.PageID, path); err != nil {
			return CreateResult{}, err
		}
		if err := s.Notion.MarkInvoiceStatus(ctx, draft.PageID, "Sent"); err != nil {
			return CreateResult{}, err
		}
	}

	res := CreateResult{
		PageID:         draft.PageID,
		InvoiceNumber:  draft.InvoiceNumber,
		PDFPath:        path,
		PDFSizeBytes:   size,
		IdempotencyKey: key,
		Totals:         totals,
		Meta:           map[string]any{},
	}
	if usedFallback {
		res.Meta["warning"] = "DEPENDENCY_FALLBACK_USED"
	}
	if strings.TrimSpace(source) == "local" {
		res.Meta["source"] = "local"
	} else {
		res.Meta["source"] = fmt.Sprintf("%s", source)
	}

	store[key] = idempotencyRecord{Key: key, Result: res, StoredAt: s.Clock.Now()}
	if err := s.saveIdempotency(store); err != nil {
		return CreateResult{}, perr.Wrap(perr.KindInternal, "failed to save idempotency store", err)
	}
	return res, nil
}
