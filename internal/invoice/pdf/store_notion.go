package pdf

import (
	"context"

	perr "biz/internal/platform/errors"
)

type Uploader interface {
	UploadInvoicePDF(ctx context.Context, pageID, path string) error
}

type NotionStore struct {
	Client Uploader
}

func (s NotionStore) Store(ctx context.Context, pageID, pdfPath string) error {
	if s.Client == nil {
		return perr.New(perr.KindDependencyUnavailable, "notion upload is not configured")
	}
	return s.Client.UploadInvoicePDF(ctx, pageID, pdfPath)
}
