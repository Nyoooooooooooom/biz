package policy

import (
	"regexp"
	"strings"

	"biz/internal/platform/config"
	perr "biz/internal/platform/errors"
)

type AgentAuthorizer struct {
	enabled                   bool
	allowedCommands           map[string]struct{}
	invoiceIDRegex            *regexp.Regexp
	maxListLimit              int
	recordsAllowedCollections map[string]struct{}
	recordsAllowedProperties  map[string]struct{}
}

type Request struct {
	Actor      string
	Command    string
	InvoiceID  string
	ListLimit  int
	Collection string
	Properties []string
}

func NewAgentAuthorizer(cfg config.AgentPolicy) (*AgentAuthorizer, error) {
	a := &AgentAuthorizer{
		enabled:                   cfg.Enabled,
		allowedCommands:           map[string]struct{}{},
		maxListLimit:              cfg.MaxListLimit,
		recordsAllowedCollections: map[string]struct{}{},
		recordsAllowedProperties:  map[string]struct{}{},
	}
	for _, c := range cfg.AllowedCommands {
		c = strings.TrimSpace(strings.ToLower(c))
		if c == "" {
			continue
		}
		a.allowedCommands[c] = struct{}{}
	}
	for _, c := range cfg.RecordsAllowedCollections {
		c = strings.TrimSpace(strings.ToLower(c))
		if c == "" {
			continue
		}
		a.recordsAllowedCollections[c] = struct{}{}
	}
	for _, p := range cfg.RecordsAllowedProperties {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		a.recordsAllowedProperties[p] = struct{}{}
	}
	if strings.TrimSpace(cfg.InvoiceIDRegex) != "" {
		r, err := regexp.Compile(cfg.InvoiceIDRegex)
		if err != nil {
			return nil, perr.Wrap(perr.KindValidation, "invalid agent_policy.invoice_id_regex", err)
		}
		a.invoiceIDRegex = r
	}
	return a, nil
}

func (a *AgentAuthorizer) Enforce(r Request) error {
	if a == nil || !a.enabled {
		return nil
	}
	if !strings.EqualFold(strings.TrimSpace(r.Actor), "agent") {
		return nil
	}
	cmd := strings.TrimSpace(strings.ToLower(r.Command))
	if _, ok := a.allowedCommands[cmd]; !ok {
		return perr.New(perr.KindValidation, "agent policy denied command: "+cmd)
	}
	if id := strings.TrimSpace(r.InvoiceID); id != "" && a.invoiceIDRegex != nil {
		if !a.invoiceIDRegex.MatchString(id) {
			return perr.New(perr.KindValidation, "agent policy denied invoice_id by regex")
		}
	}
	if (cmd == "invoice.list" || cmd == "records.list") && a.maxListLimit > 0 && r.ListLimit > a.maxListLimit {
		return perr.New(perr.KindValidation, "agent policy denied list limit above max_list_limit")
	}
	if cmd == "records.create" || cmd == "records.update" || cmd == "records.archive" {
		collection := strings.TrimSpace(strings.ToLower(r.Collection))
		if collection == "" {
			return perr.New(perr.KindValidation, "agent policy denied records mutation without collection")
		}
		if len(a.recordsAllowedCollections) == 0 {
			return perr.New(perr.KindValidation, "agent policy denied records mutation: records_allowed_collections is empty")
		}
		if _, ok := a.recordsAllowedCollections[collection]; !ok {
			return perr.New(perr.KindValidation, "agent policy denied records collection")
		}
	}
	if cmd == "records.create" || cmd == "records.update" {
		if len(a.recordsAllowedProperties) == 0 {
			return perr.New(perr.KindValidation, "agent policy denied records mutation: records_allowed_properties is empty")
		}
		for _, prop := range r.Properties {
			prop = strings.TrimSpace(prop)
			if prop == "" {
				continue
			}
			if _, ok := a.recordsAllowedProperties[prop]; !ok {
				return perr.New(perr.KindValidation, "agent policy denied records property: "+prop)
			}
		}
	}
	return nil
}
