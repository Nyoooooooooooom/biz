package tax

import "context"

type TaxInput struct {
	Region   string
	Subtotal float64
	Currency string
}

type TaxOutput struct {
	Rate   float64
	Amount float64
	Region string
}

type Fragment interface {
	Apply(ctx context.Context, in TaxInput) (TaxOutput, error)
}
