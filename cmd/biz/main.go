package main

import (
	"os"
	"strings"

	"biz/internal/command"
	invoicemodule "biz/internal/modules/invoice"
	recordsmodule "biz/internal/modules/records"
	taxmodule "biz/internal/modules/tax"
	"biz/internal/platform/config"
)

func main() {
	configPath, profile := parseBootstrapArgs(os.Args[1:])
	modules := defaultModules()
	cfg, err := config.Load(configPath, profile)
	if err == nil {
		modules = modulesFromEnabled(cfg.Modules.Enabled)
	}
	os.Exit(command.Execute(modules))
}

func parseBootstrapArgs(args []string) (configPath, profile string) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--config" && i+1 < len(args):
			configPath = args[i+1]
			i++
		case a == "--profile" && i+1 < len(args):
			profile = args[i+1]
			i++
		case len(a) > 9 && a[:9] == "--config=":
			configPath = a[9:]
		case len(a) > 10 && a[:10] == "--profile=":
			profile = a[10:]
		}
	}
	return configPath, profile
}

func defaultModules() []command.Module {
	return []command.Module{
		taxmodule.New(),
		invoicemodule.New(),
	}
}

func modulesFromEnabled(enabled []string) []command.Module {
	seen := map[string]bool{}
	var add func([]command.Module, string) []command.Module
	add = func(out []command.Module, name string) []command.Module {
		if seen[name] {
			return out
		}
		seen[name] = true
		switch name {
		case "tax":
			return append(out, taxmodule.New())
		case "invoice":
			// invoice depends on tax runtime wiring.
			out = add(out, "tax")
			return append(out, invoicemodule.New())
		case "records":
			return append(out, recordsmodule.New())
		default:
			return out
		}
	}

	var out []command.Module
	for _, name := range enabled {
		normalized := strings.ToLower(strings.TrimSpace(name))
		if normalized == "" {
			continue
		}
		out = add(out, normalized)
	}
	if len(out) == 0 {
		return defaultModules()
	}
	return out
}
