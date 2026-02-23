package invoice

import (
	"context"

	perr "biz/internal/platform/errors"
)

func (s *Service) workflowList(ctx context.Context, req ListRequest) (ListResult, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	status := normalizeStatus(req.Status)
	source := s.resolveSource("")
	if source == "local" {
		if s.FallbackFile == "" {
			return ListResult{}, perr.New(perr.KindValidation, "no local fallback file configured")
		}
		items, next, err := s.LocalData.ListByStatus(ctx, s.FallbackFile, status, limit, req.Cursor)
		if err != nil {
			return ListResult{}, err
		}
		return ListResult{Items: items, NextCursor: next}, nil
	}
	items, next, err := s.Notion.ListInvoices(ctx, status, limit, req.Cursor)
	if err == nil {
		return ListResult{Items: items, NextCursor: next}, nil
	}
	if s.FallbackFile == "" {
		return ListResult{}, err
	}
	items, next, ferr := s.LocalData.ListByStatus(ctx, s.FallbackFile, status, limit, req.Cursor)
	if ferr != nil {
		return ListResult{}, err
	}
	return ListResult{Items: items, NextCursor: next}, nil
}
