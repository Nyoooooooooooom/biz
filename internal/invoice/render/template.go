package render

import (
	"bytes"
	"context"
	"html/template"
	"math"
	"path/filepath"

	"biz/internal/invoice"
	perr "biz/internal/platform/errors"
)

type Template struct {
	path string
}

func NewTemplate(path string) *Template {
	return &Template{path: path}
}

func (t *Template) RenderInvoiceHTML(_ context.Context, doc invoice.InvoiceDocument) ([]byte, error) {
	tpl, err := template.New("invoice").Funcs(template.FuncMap{
		"lineTotal": func(q, r, d float64) float64 {
			return math.Round((q*r-d)*100) / 100
		},
	}).ParseFiles(t.path)
	if err != nil {
		return nil, perr.Wrap(perr.KindInternal, "failed to parse invoice template", err)
	}
	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, filepath.Base(t.path), doc); err != nil {
		return nil, perr.Wrap(perr.KindInternal, "failed to execute invoice template", err)
	}
	return buf.Bytes(), nil
}
