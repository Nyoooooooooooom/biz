package policy

import (
	"testing"

	"biz/internal/platform/config"
)

func TestAgentAuthorizer(t *testing.T) {
	a, err := NewAgentAuthorizer(config.AgentPolicy{
		Enabled:         true,
		AllowedCommands: []string{"invoice.list", "invoice.preview"},
		InvoiceIDRegex:  "^[a-zA-Z0-9-]{8,64}$",
		MaxListLimit:    10,
	})
	if err != nil {
		t.Fatalf("new authorizer: %v", err)
	}
	if err := a.Enforce(Request{Actor: "agent", Command: "invoice.list", ListLimit: 11}); err == nil {
		t.Fatal("expected list limit denial")
	}
	if err := a.Enforce(Request{Actor: "agent", Command: "invoice.create", InvoiceID: "abcde12345"}); err == nil {
		t.Fatal("expected command denial")
	}
	if err := a.Enforce(Request{Actor: "agent", Command: "invoice.preview", InvoiceID: "bad"}); err == nil {
		t.Fatal("expected invoice id regex denial")
	}
	if err := a.Enforce(Request{Actor: "human", Command: "invoice.create", InvoiceID: "bad"}); err != nil {
		t.Fatalf("human should bypass agent policy: %v", err)
	}
}

func TestAgentAuthorizerRecordsWriteAllowlist(t *testing.T) {
	a, err := NewAgentAuthorizer(config.AgentPolicy{
		Enabled:                   true,
		AllowedCommands:           []string{"records.create", "records.update", "records.archive"},
		RecordsAllowedCollections: []string{"invoices"},
		RecordsAllowedProperties:  []string{"Status", "Notes"},
	})
	if err != nil {
		t.Fatalf("new authorizer: %v", err)
	}

	if err := a.Enforce(Request{
		Actor:      "agent",
		Command:    "records.create",
		Collection: "invoices",
		Properties: []string{"Status"},
	}); err != nil {
		t.Fatalf("expected records.create allowed, got: %v", err)
	}
	if err := a.Enforce(Request{
		Actor:      "agent",
		Command:    "records.update",
		Collection: "invoices",
		Properties: []string{"SecretField"},
	}); err == nil {
		t.Fatal("expected denied property")
	}
	if err := a.Enforce(Request{
		Actor:      "agent",
		Command:    "records.archive",
		Collection: "clients",
	}); err == nil {
		t.Fatal("expected denied collection")
	}
}

func TestAgentAuthorizerRecordsListLimit(t *testing.T) {
	a, err := NewAgentAuthorizer(config.AgentPolicy{
		Enabled:         true,
		AllowedCommands: []string{"records.list"},
		MaxListLimit:    10,
	})
	if err != nil {
		t.Fatalf("new authorizer: %v", err)
	}
	if err := a.Enforce(Request{Actor: "agent", Command: "records.list", ListLimit: 11}); err == nil {
		t.Fatal("expected records.list limit denial")
	}
	if err := a.Enforce(Request{Actor: "agent", Command: "records.list", ListLimit: 5}); err != nil {
		t.Fatalf("expected records.list allowed: %v", err)
	}
}
