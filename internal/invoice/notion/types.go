package notion

import "time"

type rawPage struct {
	ID             string         `json:"id"`
	LastEditedTime time.Time      `json:"last_edited_time"`
	Properties     map[string]any `json:"properties"`
}
