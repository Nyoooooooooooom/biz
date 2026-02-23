package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	perr "biz/internal/platform/errors"
	"biz/internal/records"
)

func (c *Client) List(ctx context.Context, dbID string, limit int, cursor string) ([]records.Record, string, error) {
	if strings.TrimSpace(dbID) == "" {
		return nil, "", perr.New(perr.KindValidation, "notion db id is required")
	}
	payload := map[string]any{"page_size": limit}
	if strings.TrimSpace(cursor) != "" {
		payload["start_cursor"] = cursor
	}
	body, _ := json.Marshal(payload)
	b, err := c.do(ctx, "POST", fmt.Sprintf("%s/databases/%s/query", c.BaseURL, dbID), body)
	if err != nil {
		return nil, "", err
	}
	var qr struct {
		Results    []rawPage `json:"results"`
		NextCursor string    `json:"next_cursor"`
	}
	if err := json.Unmarshal(b, &qr); err != nil {
		return nil, "", perr.Wrap(perr.KindDependencyUnavailable, "invalid notion list response", err)
	}
	items := make([]records.Record, 0, len(qr.Results))
	for _, p := range qr.Results {
		items = append(items, mapPage(p))
	}
	return items, qr.NextCursor, nil
}

func (c *Client) Get(ctx context.Context, id string) (records.Record, error) {
	if strings.TrimSpace(id) == "" {
		return records.Record{}, perr.New(perr.KindValidation, "notion page id is required")
	}
	b, err := c.do(ctx, "GET", fmt.Sprintf("%s/pages/%s", c.BaseURL, id), nil)
	if err != nil {
		return records.Record{}, err
	}
	var p rawPage
	if err := json.Unmarshal(b, &p); err != nil {
		return records.Record{}, perr.Wrap(perr.KindDependencyUnavailable, "invalid notion page response", err)
	}
	return mapPage(p), nil
}

func (c *Client) Create(ctx context.Context, dbID string, properties map[string]any) (records.Record, error) {
	if strings.TrimSpace(dbID) == "" {
		return records.Record{}, perr.New(perr.KindValidation, "notion db id is required")
	}
	payload := map[string]any{
		"parent": map[string]any{
			"database_id": dbID,
		},
		"properties": properties,
	}
	body, _ := json.Marshal(payload)
	b, err := c.do(ctx, "POST", fmt.Sprintf("%s/pages", c.BaseURL), body)
	if err != nil {
		return records.Record{}, err
	}
	var p rawPage
	if err := json.Unmarshal(b, &p); err != nil {
		return records.Record{}, perr.Wrap(perr.KindDependencyUnavailable, "invalid notion create response", err)
	}
	return mapPage(p), nil
}

func (c *Client) Update(ctx context.Context, id string, properties map[string]any) (records.Record, error) {
	if strings.TrimSpace(id) == "" {
		return records.Record{}, perr.New(perr.KindValidation, "notion page id is required")
	}
	payload := map[string]any{"properties": properties}
	body, _ := json.Marshal(payload)
	b, err := c.do(ctx, "PATCH", fmt.Sprintf("%s/pages/%s", c.BaseURL, id), body)
	if err != nil {
		return records.Record{}, err
	}
	var p rawPage
	if err := json.Unmarshal(b, &p); err != nil {
		return records.Record{}, perr.Wrap(perr.KindDependencyUnavailable, "invalid notion update response", err)
	}
	return mapPage(p), nil
}

func (c *Client) Archive(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return perr.New(perr.KindValidation, "notion page id is required")
	}
	payload := map[string]any{"archived": true}
	body, _ := json.Marshal(payload)
	_, err := c.do(ctx, "PATCH", fmt.Sprintf("%s/pages/%s", c.BaseURL, id), body)
	return err
}

func (c *Client) Schema(ctx context.Context, dbID string) (records.Schema, error) {
	if strings.TrimSpace(dbID) == "" {
		return records.Schema{}, perr.New(perr.KindValidation, "notion db id is required")
	}
	b, err := c.do(ctx, "GET", fmt.Sprintf("%s/databases/%s", c.BaseURL, dbID), nil)
	if err != nil {
		return records.Schema{}, err
	}
	var db struct {
		ID         string                    `json:"id"`
		Properties map[string]map[string]any `json:"properties"`
	}
	if err := json.Unmarshal(b, &db); err != nil {
		return records.Schema{}, perr.Wrap(perr.KindDependencyUnavailable, "invalid notion schema response", err)
	}
	out := records.Schema{DBID: db.ID, Properties: map[string]string{}}
	if out.DBID == "" {
		out.DBID = dbID
	}
	for name, raw := range db.Properties {
		t, _ := raw["type"].(string)
		out.Properties[name] = t
	}
	return out, nil
}
