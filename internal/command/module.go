package command

import "biz/internal/app"

import "github.com/spf13/cobra"

// Module is a pluggable command fragment mounted under the root CLI.
type Module interface {
	Name() string
	Build(*State) *cobra.Command
}

// RuntimeInitializer optionally wires module-specific runtime dependencies.
type RuntimeInitializer interface {
	InitRuntime(*app.Runtime) error
}
