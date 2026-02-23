package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

type FlagType string

const (
	FlagString FlagType = "string"
	FlagBool   FlagType = "bool"
	FlagInt    FlagType = "int"
)

type FlagSpec struct {
	Type      FlagType
	Name      string
	Shorthand string
	Usage     string
	Required  bool
	Target    any
	Default   any
}

type CommandSpec struct {
	Use      string
	Aliases  []string
	Short    string
	Args     cobra.PositionalArgs
	RunE     func(cmd *cobra.Command, args []string) error
	Flags    []FlagSpec
	Commands []CommandSpec
}

func BuildCommand(spec CommandSpec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     spec.Use,
		Aliases: spec.Aliases,
		Short:   spec.Short,
		Args:    spec.Args,
		RunE:    spec.RunE,
	}
	for _, f := range spec.Flags {
		applyFlag(cmd, f)
	}
	for _, sub := range spec.Commands {
		cmd.AddCommand(BuildCommand(sub))
	}
	return cmd
}

func applyFlag(cmd *cobra.Command, f FlagSpec) {
	if f.Shorthand != "" && len(f.Shorthand) != 1 {
		panic(fmt.Sprintf("invalid shorthand for flag %s", f.Name))
	}
	sw := cmd.Flags()
	switch ptr := f.Target.(type) {
	case *string:
		dv, _ := f.Default.(string)
		if f.Shorthand == "" {
			sw.StringVar(ptr, f.Name, dv, f.Usage)
		} else {
			sw.StringVarP(ptr, f.Name, f.Shorthand, dv, f.Usage)
		}
	case *bool:
		dv, _ := f.Default.(bool)
		if f.Shorthand == "" {
			sw.BoolVar(ptr, f.Name, dv, f.Usage)
		} else {
			sw.BoolVarP(ptr, f.Name, f.Shorthand, dv, f.Usage)
		}
	case *int:
		dv, _ := f.Default.(int)
		if f.Shorthand == "" {
			sw.IntVar(ptr, f.Name, dv, f.Usage)
		} else {
			sw.IntVarP(ptr, f.Name, f.Shorthand, dv, f.Usage)
		}
	default:
		panic(fmt.Sprintf("unsupported flag target type for %s", f.Name))
	}
	if f.Required {
		_ = cmd.MarkFlagRequired(f.Name)
	}
}
