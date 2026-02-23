package records

import (
	"biz/internal/app"
	"biz/internal/command"
	recordsdomain "biz/internal/records"
	recordsnotion "biz/internal/records/notion"
	"github.com/spf13/cobra"
)

type Module struct{}

func New() command.Module { return Module{} }

func (Module) Name() string { return "records" }

func (Module) InitRuntime(rt *app.Runtime) error {
	cfg := rt.Config
	rt.Records = recordsdomain.Service{
		Reader: recordsnotion.New(
			cfg.Notion.Token,
			cfg.Notion.BaseURL,
			cfg.Notion.ReadTimeout,
			cfg.Notion.RetryCount,
			cfg.Notion.RetryBackoffMS,
		),
	}
	return nil
}

func (m Module) Build(s *command.State) *cobra.Command {
	return buildRecordsCommand(s)
}
