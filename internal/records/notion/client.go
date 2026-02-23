package notion

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	perr "biz/internal/platform/errors"
)

const notionVersion = "2022-06-28"

type Client struct {
	HTTP       *http.Client
	Token      string
	BaseURL    string
	RetryCount int
	BackoffMS  int
	rand       *rand.Rand
}

func New(token, baseURL string, timeout time.Duration, retryCount, backoffMS int) *Client {
	if baseURL == "" {
		baseURL = "https://api.notion.com/v1"
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	if retryCount <= 0 {
		retryCount = 3
	}
	if backoffMS <= 0 {
		backoffMS = 200
	}
	return &Client{
		HTTP:       &http.Client{Timeout: timeout},
		Token:      token,
		BaseURL:    strings.TrimRight(baseURL, "/"),
		RetryCount: retryCount,
		BackoffMS:  backoffMS,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *Client) do(ctx context.Context, method, url string, body []byte) ([]byte, error) {
	if strings.TrimSpace(c.Token) == "" {
		return nil, perr.New(perr.KindDependencyUnavailable, "notion not configured")
	}
	for i := 0; i < c.RetryCount; i++ {
		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
		if err != nil {
			return nil, perr.Wrap(perr.KindInternal, "failed to build notion request", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.Token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Notion-Version", notionVersion)
		resp, err := c.HTTP.Do(req)
		if err != nil {
			if i < c.RetryCount-1 {
				time.Sleep(c.backoff(i))
				continue
			}
			return nil, perr.Wrap(perr.KindDependencyUnavailable, "notion request failed", err)
		}
		b, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			if i < c.RetryCount-1 {
				time.Sleep(c.backoff(i))
				continue
			}
			return nil, perr.Wrap(perr.KindDependencyUnavailable, "failed reading notion response", readErr)
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return b, nil
		}
		if i < c.RetryCount-1 && (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500) {
			time.Sleep(c.backoff(i))
			continue
		}
		if resp.StatusCode == http.StatusNotFound {
			return nil, perr.New(perr.KindNotFound, "notion resource not found")
		}
		return nil, perr.New(perr.KindDependencyUnavailable, fmt.Sprintf("notion error: %s", resp.Status))
	}
	return nil, perr.New(perr.KindDependencyUnavailable, "notion unavailable")
}

func (c *Client) backoff(attempt int) time.Duration {
	base := time.Duration(c.BackoffMS*(1<<attempt)) * time.Millisecond
	jitter := time.Duration(c.rand.Intn(c.BackoffMS+1)) * time.Millisecond
	return base + jitter
}
