package main

import "testing"

func TestParseBootstrapArgs(t *testing.T) {
	configPath, profile := parseBootstrapArgs([]string{
		"--config", "a.yaml",
		"--profile=prod",
	})
	if configPath != "a.yaml" {
		t.Fatalf("unexpected config path: %q", configPath)
	}
	if profile != "prod" {
		t.Fatalf("unexpected profile: %q", profile)
	}
}

func TestModulesFromEnabledIncludesInvoiceDependency(t *testing.T) {
	modules := modulesFromEnabled([]string{"invoice"})
	if len(modules) != 2 {
		t.Fatalf("expected tax+invoice modules, got %d", len(modules))
	}
	if modules[0].Name() != "tax" {
		t.Fatalf("expected first module tax, got %s", modules[0].Name())
	}
	if modules[1].Name() != "invoice" {
		t.Fatalf("expected second module invoice, got %s", modules[1].Name())
	}
}

func TestModulesFromEnabledMixed(t *testing.T) {
	modules := modulesFromEnabled([]string{"records", "invoice", "records"})
	got := make([]string, 0, len(modules))
	for _, m := range modules {
		got = append(got, m.Name())
	}
	want := []string{"records", "tax", "invoice"}
	if len(got) != len(want) {
		t.Fatalf("unexpected module count: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected module order: got=%v want=%v", got, want)
		}
	}
}
