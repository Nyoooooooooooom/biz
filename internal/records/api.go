package records

import (
	"context"
	"time"
)

type Record struct {
	ID             string         `json:"id"`
	LastEditedTime time.Time      `json:"last_edited_time"`
	Properties     map[string]any `json:"properties"`
}

type ListRequest struct {
	Collection string `json:"collection"`
	DBID       string `json:"db_id"`
	Limit      int    `json:"limit"`
	Cursor     string `json:"cursor,omitempty"`
}

type ListResult struct {
	Items      []Record `json:"items"`
	NextCursor string   `json:"next_cursor,omitempty"`
}

type GetRequest struct {
	ID string `json:"id"`
}

type CreateRequest struct {
	Collection string         `json:"collection"`
	DBID       string         `json:"db_id"`
	Properties map[string]any `json:"properties"`
}

type UpdateRequest struct {
	ID         string         `json:"id"`
	Collection string         `json:"collection,omitempty"`
	Properties map[string]any `json:"properties"`
}

type ArchiveRequest struct {
	ID         string `json:"id"`
	Collection string `json:"collection,omitempty"`
}

type SchemaRequest struct {
	Collection string `json:"collection"`
	DBID       string `json:"db_id"`
}

type Schema struct {
	DBID       string            `json:"db_id"`
	Properties map[string]string `json:"properties"`
}

type Fragment interface {
	List(ctx context.Context, req ListRequest) (ListResult, error)
	Get(ctx context.Context, req GetRequest) (Record, error)
	Create(ctx context.Context, req CreateRequest) (Record, error)
	Update(ctx context.Context, req UpdateRequest) (Record, error)
	Archive(ctx context.Context, req ArchiveRequest) error
	Schema(ctx context.Context, req SchemaRequest) (Schema, error)
}
