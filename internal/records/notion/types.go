package notion

import (
	"time"

	"biz/internal/records"
)

type rawPage struct {
	ID             string         `json:"id"`
	LastEditedTime time.Time      `json:"last_edited_time"`
	Properties     map[string]any `json:"properties"`
}

func mapPage(p rawPage) records.Record {
	return records.Record{
		ID:             p.ID,
		LastEditedTime: p.LastEditedTime,
		Properties:     p.Properties,
	}
}
