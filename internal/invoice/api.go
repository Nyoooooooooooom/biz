package invoice

import "context"

type Fragment interface {
	Create(ctx context.Context, req CreateRequest) (CreateResult, error)
	List(ctx context.Context, req ListRequest) (ListResult, error)
	Preview(ctx context.Context, req PreviewRequest) (PreviewResult, error)
	Validate(ctx context.Context, in InvoiceDraft) error
}
