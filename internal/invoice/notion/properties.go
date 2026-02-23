package notion

import (
	"strconv"
	"strings"
	"time"
)

func relationIDs(props map[string]any, key string) []string {
	p := property(props, key)
	raw, ok := p["relation"].([]any)
	if !ok {
		return nil
	}
	ids := make([]string, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		id, _ := m["id"].(string)
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func property(props map[string]any, key string) map[string]any {
	if props == nil {
		return nil
	}
	raw, ok := props[key]
	if !ok {
		return nil
	}
	m, _ := raw.(map[string]any)
	return m
}

func propText(props map[string]any, key string) string {
	p := property(props, key)
	if p == nil {
		return ""
	}
	if selectObj, ok := p["select"].(map[string]any); ok {
		if name, ok := selectObj["name"].(string); ok {
			return strings.TrimSpace(name)
		}
	}
	if statusObj, ok := p["status"].(map[string]any); ok {
		if name, ok := statusObj["name"].(string); ok {
			return strings.TrimSpace(name)
		}
	}
	if title, ok := p["title"].([]any); ok {
		if s := textArray(title); s != "" {
			return s
		}
	}
	if rich, ok := p["rich_text"].([]any); ok {
		if s := textArray(rich); s != "" {
			return s
		}
	}
	if plain, ok := p["plain_text"].(string); ok {
		return strings.TrimSpace(plain)
	}
	if n, ok := p["name"].(string); ok {
		return strings.TrimSpace(n)
	}
	if v, ok := p["value"].(string); ok {
		return strings.TrimSpace(v)
	}
	return ""
}

func textArray(arr []any) string {
	parts := make([]string, 0, len(arr))
	for _, item := range arr {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if p, ok := m["plain_text"].(string); ok && strings.TrimSpace(p) != "" {
			parts = append(parts, strings.TrimSpace(p))
			continue
		}
		if t, ok := m["text"].(map[string]any); ok {
			if c, ok := t["content"].(string); ok && strings.TrimSpace(c) != "" {
				parts = append(parts, strings.TrimSpace(c))
			}
		}
	}
	return strings.Join(parts, "")
}

func propDate(props map[string]any, key string) time.Time {
	p := property(props, key)
	if p == nil {
		return time.Time{}
	}
	if d, ok := p["date"].(map[string]any); ok {
		if s, ok := d["start"].(string); ok {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				return t
			}
			if t, err := time.Parse("2006-01-02", s); err == nil {
				return t
			}
		}
	}
	if s, ok := p["value"].(string); ok {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t
		}
		if t, err := time.Parse("2006-01-02", s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func firstText(props map[string]any, keys ...string) string {
	for _, k := range keys {
		if v := propText(props, k); v != "" {
			return v
		}
	}
	return ""
}

func firstDate(props map[string]any, keys ...string) time.Time {
	for _, k := range keys {
		if v := propDate(props, k); !v.IsZero() {
			return v
		}
	}
	return time.Time{}
}

func firstNumber(props map[string]any, keys ...string) float64 {
	for _, k := range keys {
		if p := property(props, k); p != nil {
			if n, ok := p["number"].(float64); ok {
				return n
			}
			if f, ok := p["formula"].(map[string]any); ok {
				if n, ok := f["number"].(float64); ok {
					return n
				}
				if s, ok := f["string"].(string); ok && s != "" {
					if n, err := strconv.ParseFloat(strings.TrimSpace(s), 64); err == nil {
						return n
					}
				}
			}
			if r, ok := p["rollup"].(map[string]any); ok {
				if n, ok := r["number"].(float64); ok {
					return n
				}
				if arr, ok := r["array"].([]any); ok {
					total := 0.0
					found := false
					for _, item := range arr {
						m, ok := item.(map[string]any)
						if !ok {
							continue
						}
						if n, ok := m["number"].(float64); ok {
							total += n
							found = true
						}
					}
					if found {
						return total
					}
				}
			}
			if n, ok := p["value"].(float64); ok {
				return n
			}
		}
	}
	return 0
}
