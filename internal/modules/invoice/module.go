package invoice

import (
	"biz/internal/command"
	"github.com/spf13/cobra"
)

type Module struct{}

func New() command.Module { return Module{} }

func (Module) Name() string { return "invoice" }

func (m Module) Build(s *command.State) *cobra.Command {
	return buildInvoiceCommand(s)
}
