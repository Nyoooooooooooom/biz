package records

import (
	"context"
	"strings"

	perr "biz/internal/platform/errors"
)

type Reader interface {
	List(ctx context.Context, dbID string, limit int, cursor string) ([]Record, string, error)
	Get(ctx context.Context, id string) (Record, error)
	Create(ctx context.Context, dbID string, properties map[string]any) (Record, error)
	Update(ctx context.Context, id string, properties map[string]any) (Record, error)
	Archive(ctx context.Context, id string) error
	Schema(ctx context.Context, dbID string) (Schema, error)
}

type Service struct {
	Reader Reader
}

func (s Service) List(ctx context.Context, req ListRequest) (ListResult, error) {
	if s.Reader == nil {
		return ListResult{}, perr.New(perr.KindDependencyUnavailable, "records reader is not configured")
	}
	if strings.TrimSpace(req.DBID) == "" {
		return ListResult{}, perr.New(perr.KindValidation, "db_id is required")
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	items, next, err := s.Reader.List(ctx, req.DBID, limit, req.Cursor)
	if err != nil {
		return ListResult{}, err
	}
	return ListResult{Items: items, NextCursor: next}, nil
}

func (s Service) Get(ctx context.Context, req GetRequest) (Record, error) {
	if s.Reader == nil {
		return Record{}, perr.New(perr.KindDependencyUnavailable, "records reader is not configured")
	}
	if strings.TrimSpace(req.ID) == "" {
		return Record{}, perr.New(perr.KindValidation, "record id is required")
	}
	return s.Reader.Get(ctx, req.ID)
}

func (s Service) Create(ctx context.Context, req CreateRequest) (Record, error) {
	if s.Reader == nil {
		return Record{}, perr.New(perr.KindDependencyUnavailable, "records reader is not configured")
	}
	if strings.TrimSpace(req.DBID) == "" {
		return Record{}, perr.New(perr.KindValidation, "db_id is required")
	}
	if len(req.Properties) == 0 {
		return Record{}, perr.New(perr.KindValidation, "properties are required")
	}
	return s.Reader.Create(ctx, req.DBID, req.Properties)
}

func (s Service) Update(ctx context.Context, req UpdateRequest) (Record, error) {
	if s.Reader == nil {
		return Record{}, perr.New(perr.KindDependencyUnavailable, "records reader is not configured")
	}
	if strings.TrimSpace(req.ID) == "" {
		return Record{}, perr.New(perr.KindValidation, "record id is required")
	}
	if len(req.Properties) == 0 {
		return Record{}, perr.New(perr.KindValidation, "properties are required")
	}
	return s.Reader.Update(ctx, req.ID, req.Properties)
}

func (s Service) Archive(ctx context.Context, req ArchiveRequest) error {
	if s.Reader == nil {
		return perr.New(perr.KindDependencyUnavailable, "records reader is not configured")
	}
	if strings.TrimSpace(req.ID) == "" {
		return perr.New(perr.KindValidation, "record id is required")
	}
	return s.Reader.Archive(ctx, req.ID)
}

func (s Service) Schema(ctx context.Context, req SchemaRequest) (Schema, error) {
	if s.Reader == nil {
		return Schema{}, perr.New(perr.KindDependencyUnavailable, "records reader is not configured")
	}
	if strings.TrimSpace(req.DBID) == "" {
		return Schema{}, perr.New(perr.KindValidation, "db_id is required")
	}
	return s.Reader.Schema(ctx, req.DBID)
}
