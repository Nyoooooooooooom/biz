package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"biz/internal/app"
	"biz/internal/audit"
	perr "biz/internal/platform/errors"
	"biz/internal/platform/id"
	"biz/internal/platform/output"
	"github.com/spf13/cobra"
)

type State struct {
	Runtime    *app.Runtime
	ConfigPath string
	Profile    string
	Actor      string
	JSONOut    bool
	TraceID    string
	Command    string
	Args       []string
}

func Execute(modules []Module) int {
	s := &State{TraceID: id.TraceID()}
	root := newRootCmd(s, modules)
	err := root.ExecuteContext(context.Background())
	exitCode := 0
	resultCode := "OK"
	errMessage := ""
	if err != nil {
		exitCode = handleErr(s, err)
		resultCode = string(perr.KindOf(err))
		if resultCode == "" {
			resultCode = "INTERNAL_ERROR"
		}
		errMessage = err.Error()
	}
	if aerr := writeAudit(s, exitCode, resultCode, errMessage); aerr != nil {
		fmt.Fprintln(os.Stderr, "audit log write failed:", aerr)
		if s != nil && s.Runtime != nil && s.Runtime.Config.Audit.Strict && exitCode == 0 {
			return 4
		}
	}
	return exitCode
}

func newRootCmd(s *State, modules []Module) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "biz",
		Short:         "Composable business automation CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return initializeStateRuntime(s, cmd, modules)
		},
	}
	cmd.PersistentFlags().StringVar(&s.ConfigPath, "config", "", "Path to config file")
	cmd.PersistentFlags().StringVar(&s.Profile, "profile", "", "Config profile override")
	cmd.PersistentFlags().StringVar(&s.Actor, "actor", "human", "Invocation actor (human|agent)")
	cmd.PersistentFlags().BoolVar(&s.JSONOut, "json", false, "Emit machine-readable JSON output")
	cmd.PersistentFlags().StringVar(&s.TraceID, "trace-id", s.TraceID, "Trace identifier")

	registerModules(cmd, s, modules)

	cmd.AddCommand(&cobra.Command{
		Use:   "doctor",
		Short: "Validate runtime configuration and dependencies",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDoctor(s)
		},
	})

	return cmd
}

func initializeStateRuntime(s *State, cmd *cobra.Command, modules []Module) error {
	rt, err := app.Build(s.ConfigPath, s.Profile)
	if err != nil {
		return err
	}
	if err := initModuleRuntime(rt, modules); err != nil {
		return err
	}
	s.Runtime = rt
	s.Command = cmd.CommandPath()
	s.Args = os.Args[1:]
	configureProductionLogger(s)
	return nil
}

func initModuleRuntime(rt *app.Runtime, modules []Module) error {
	for _, m := range modules {
		if initializer, ok := m.(RuntimeInitializer); ok {
			if err := initializer.InitRuntime(rt); err != nil {
				return err
			}
		}
	}
	return nil
}

func configureProductionLogger(s *State) {
	if s == nil || s.Runtime == nil {
		return
	}
	if strings.EqualFold(s.Runtime.Config.Profile, "prod") {
		s.Runtime.Logger = s.Runtime.Logger.With()
	}
}

func registerModules(root *cobra.Command, s *State, modules []Module) {
	for _, m := range modules {
		if m == nil {
			continue
		}
		if sub := m.Build(s); sub != nil {
			root.AddCommand(sub)
		}
	}
}

func runDoctor(s *State) error {
	checks := doctorChecks(s)
	if s.JSONOut {
		return EmitJSON(output.OK(s.TraceID, "doctor checks passed", checks))
	}
	fmt.Printf("profile=%s\n", s.Runtime.Config.Profile)
	fmt.Printf("actor=%s\n", s.Actor)
	fmt.Printf("notion_db_id_configured=%v\n", checks["notion_db_id"])
	fmt.Printf("template_path=%s\n", checks["template_path"])
	fmt.Printf("invoice_source=%s\n", checks["source"])
	fmt.Printf("audit_enabled=%v\n", checks["audit_enabled"])
	fmt.Printf("audit_path=%s\n", checks["audit_path"])
	return nil
}

func doctorChecks(s *State) map[string]any {
	return map[string]any{
		"profile":       s.Runtime.Config.Profile,
		"actor":         s.Actor,
		"notion_db_id":  s.Runtime.Config.Notion.InvoiceDBID != "",
		"template_path": s.Runtime.Config.Invoice.TemplatePath,
		"source":        s.Runtime.Config.Invoice.Source,
		"audit_enabled": s.Runtime.Config.Audit.Enabled,
		"audit_path":    s.Runtime.Config.Audit.Path,
	}
}

func EmitJSON[T any](env output.Envelope[T]) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(env)
}

func handleErr(s *State, err error) int {
	if s != nil && s.JSONOut {
		kind := string(perr.KindOf(err))
		code := kind
		if code == "" {
			code = "INTERNAL_ERROR"
		}
		env := output.Fail[map[string]any](s.TraceID, code, err.Error(), kind)
		_ = EmitJSON(env)
	}
	fmt.Fprintln(os.Stderr, err.Error())
	return perr.ExitCode(err)
}

func writeAudit(s *State, exitCode int, resultCode, errMessage string) error {
	if s == nil || s.Runtime == nil || s.Runtime.Auditor == nil {
		return nil
	}
	command := strings.TrimSpace(s.Command)
	if command == "" {
		command = "biz"
	}
	return s.Runtime.Auditor.Write(audit.Event{
		TraceID:      s.TraceID,
		Actor:        s.Actor,
		Command:      command,
		Args:         s.Args,
		ExitCode:     exitCode,
		ResultCode:   resultCode,
		ErrorMessage: errMessage,
	})
}
