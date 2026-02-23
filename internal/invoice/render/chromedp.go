package render

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"unicode"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/jung-kurt/gofpdf"
)

type ChromePDF struct {
	NoSandbox          bool
	DisableDevShmUsage bool
	DisableJavaScript  bool
	ExecPath           string
}

func (c ChromePDF) RenderPDF(ctx context.Context, html []byte) ([]byte, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
	)
	if c.NoSandbox {
		opts = append(opts, chromedp.Flag("no-sandbox", true))
	}
	if c.DisableDevShmUsage {
		opts = append(opts, chromedp.Flag("disable-dev-shm-usage", true))
	}
	if c.DisableJavaScript {
		opts = append(opts, chromedp.Flag("disable-javascript", true))
	}
	execPath := strings.TrimSpace(c.ExecPath)
	if execPath == "" {
		execPath = detectChromePath()
	}
	if execPath != "" {
		opts = append(opts, chromedp.ExecPath(execPath))
	}
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()
	browserCtx, bcancel := chromedp.NewContext(allocCtx)
	defer bcancel()
	dataURL := "data:text/html;base64," + base64.StdEncoding.EncodeToString(html)
	var pdf []byte
	err := chromedp.Run(browserCtx,
		chromedp.Navigate(dataURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			pdf = buf
			return nil
		}),
	)
	if err != nil {
		fallback, ferr := renderFallbackPDF(html)
		if ferr == nil {
			return fallback, nil
		}
		return nil, fmt.Errorf("chrome failed to start: %w; fallback failed: %v", err, ferr)
	}
	return pdf, nil
}

func detectChromePath() string {
	if p := os.Getenv("BIZ_CHROME_PATH"); p != "" {
		return p
	}
	candidates := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"/opt/homebrew/Caskroom/google-chrome/current/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/opt/homebrew/Caskroom/chromium/current/chrome-mac/Chromium.app/Contents/MacOS/Chromium",
		"google-chrome",
		"chromium-browser",
		"chromium",
	}
	if runtime.GOOS == "darwin" {
		candidates = append(candidates, "/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary")
		candidates = append(candidates, newestByGlob("/opt/homebrew/Caskroom/google-chrome/*/Google Chrome.app/Contents/MacOS/Google Chrome"))
		candidates = append(candidates, newestByGlob("/opt/homebrew/Caskroom/chromium/*/chrome-mac/Chromium.app/Contents/MacOS/Chromium"))
	}
	for _, c := range candidates {
		if c == "" {
			continue
		}
		if _, err := exec.LookPath(c); err == nil {
			return c
		}
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c
		}
	}
	return ""
}

func newestByGlob(pattern string) string {
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return ""
	}
	sort.Strings(matches)
	return matches[len(matches)-1]
}

func renderFallbackPDF(html []byte) ([]byte, error) {
	p := gofpdf.New("P", "mm", "A4", "")
	p.SetTitle("Invoice", false)
	p.AddPage()
	p.SetFont("Arial", "", 11)
	p.MultiCell(190, 5, htmlToText(string(html)), "", "L", false)
	var out bytes.Buffer
	if err := p.Output(&out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func htmlToText(in string) string {
	var b strings.Builder
	b.Grow(len(in))
	inTag := false
	lastSpace := false
	for _, r := range in {
		switch {
		case r == '<':
			inTag = true
			continue
		case r == '>':
			inTag = false
			if !lastSpace {
				b.WriteRune(' ')
				lastSpace = true
			}
			continue
		}
		if inTag {
			continue
		}
		if unicode.IsSpace(r) {
			if !lastSpace {
				b.WriteRune(' ')
				lastSpace = true
			}
			continue
		}
		b.WriteRune(r)
		lastSpace = false
	}
	return strings.TrimSpace(b.String())
}
