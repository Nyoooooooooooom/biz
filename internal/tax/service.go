package tax

import (
	"context"
	"math"

	perr "biz/internal/platform/errors"
)

type Service struct {
	Rates         map[string]float64
	DefaultRegion string
	Required      bool
}

func (s Service) Apply(_ context.Context, in TaxInput) (TaxOutput, error) {
	region := in.Region
	if region == "" {
		region = s.DefaultRegion
	}
	rate, ok := s.Rates[region]
	if !ok {
		if s.Required {
			return TaxOutput{}, perr.New(perr.KindValidation, "no tax rate configured for region")
		}
		rate = 0
	}
	amount := math.Round((in.Subtotal*rate)*100) / 100
	return TaxOutput{Rate: rate, Amount: amount, Region: region}, nil
}
