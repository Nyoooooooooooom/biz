package tax

import (
	"biz/internal/app"
	"biz/internal/command"
	taxsvc "biz/internal/tax"
	"github.com/spf13/cobra"
)

type Module struct{}

func New() command.Module { return Module{} }

func (Module) Name() string { return "tax" }

func (Module) InitRuntime(rt *app.Runtime) error {
	cfg := rt.Config
	rt.Tax = taxsvc.Service{
		Rates:         cfg.Tax.Rates,
		DefaultRegion: cfg.Tax.DefaultRegion,
		Required:      cfg.Tax.Required,
	}
	return nil
}

func (Module) Build(*command.State) *cobra.Command { return nil }
