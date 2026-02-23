package records

import (
	"fmt"
	"sort"
	"strings"

	"biz/internal/command"
	perr "biz/internal/platform/errors"
	"biz/internal/platform/output"
	"biz/internal/policy"
	recordsdomain "biz/internal/records"
	"github.com/spf13/cobra"
)

func buildRecordsCommand(s *command.State) *cobra.Command {
	return command.BuildCommand(command.CommandSpec{
		Use:     "records",
		Aliases: []string{"rec"},
		Short:   "Read and mutate records from configured data collections",
		Commands: []command.CommandSpec{
			buildListSpec(s),
			buildGetSpec(s),
			buildSchemaSpec(s),
			buildCreateSpec(s),
			buildUpdateSpec(s),
			buildArchiveSpec(s),
		},
	})
}

func buildListSpec(s *command.State) command.CommandSpec {
	var req recordsdomain.ListRequest
	return command.CommandSpec{
		Use:   "list <collection-or-db-id>",
		Short: "List records from a Notion collection/database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := s.Runtime.AgentAuthorizer.Enforce(policy.Request{
				Actor:     s.Actor,
				Command:   "records.list",
				ListLimit: req.Limit,
			}); err != nil {
				return err
			}
			req.Collection = args[0]
			req.DBID = resolveCollectionDBID(s, req.Collection)
			res, err := s.Runtime.Records.List(cmd.Context(), req)
			if err != nil {
				return err
			}
			if s.JSONOut {
				return command.EmitJSON(output.OK(s.TraceID, "records listed", res))
			}
			fmt.Printf("collection=%s count=%d\n", req.Collection, len(res.Items))
			for _, item := range res.Items {
				fmt.Println(item.ID)
			}
			if res.NextCursor != "" {
				fmt.Printf("next_cursor=%s\n", res.NextCursor)
			}
			return nil
		},
		Flags: []command.FlagSpec{
			{Type: command.FlagInt, Name: "limit", Usage: "Maximum results", Target: &req.Limit, Default: 20},
			{Type: command.FlagString, Name: "cursor", Usage: "Pagination cursor", Target: &req.Cursor},
		},
	}
}

func buildGetSpec(s *command.State) command.CommandSpec {
	var req recordsdomain.GetRequest
	return command.CommandSpec{
		Use:   "get <page-id>",
		Short: "Get one Notion page record by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := s.Runtime.AgentAuthorizer.Enforce(policy.Request{
				Actor:   s.Actor,
				Command: "records.get",
			}); err != nil {
				return err
			}
			req.ID = args[0]
			res, err := s.Runtime.Records.Get(cmd.Context(), req)
			if err != nil {
				return err
			}
			if s.JSONOut {
				return command.EmitJSON(output.OK(s.TraceID, "record fetched", res))
			}
			fmt.Printf("id=%s\n", res.ID)
			return nil
		},
	}
}

func buildSchemaSpec(s *command.State) command.CommandSpec {
	var req recordsdomain.SchemaRequest
	return command.CommandSpec{
		Use:   "schema <collection-or-db-id>",
		Short: "Get Notion database schema (property names/types)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req.Collection = args[0]
			req.DBID = resolveCollectionDBID(s, req.Collection)
			res, err := s.Runtime.Records.Schema(cmd.Context(), req)
			if err != nil {
				return err
			}
			if s.JSONOut {
				return command.EmitJSON(output.OK(s.TraceID, "records schema fetched", res))
			}
			fmt.Printf("db_id=%s properties=%d\n", res.DBID, len(res.Properties))
			keys := make([]string, 0, len(res.Properties))
			for k := range res.Properties {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Printf("%s: %s\n", k, res.Properties[k])
			}
			return nil
		},
	}
}

func buildCreateSpec(s *command.State) command.CommandSpec {
	var req recordsdomain.CreateRequest
	var rawData, rawFile string
	var validateSchema, dryRun bool

	return command.CommandSpec{
		Use:   "create <collection-or-db-id>",
		Short: "Create one Notion page record",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, err := parseProperties(rawData, rawFile)
			if err != nil {
				return err
			}
			req.Collection = args[0]
			req.DBID = resolveCollectionDBID(s, req.Collection)
			if validateSchema {
				if err := validatePropertiesAgainstSchema(cmd, s, req.Collection, req.DBID, props); err != nil {
					return err
				}
			}
			req.Properties = props
			if err := s.Runtime.AgentAuthorizer.Enforce(policy.Request{
				Actor:      s.Actor,
				Command:    "records.create",
				Collection: normalizeCollection(req.Collection),
				Properties: propertyKeys(req.Properties),
			}); err != nil {
				return err
			}
			if dryRun {
				payload := map[string]any{
					"action":       "create",
					"collection":   req.Collection,
					"db_id":        req.DBID,
					"propertyKeys": propertyKeys(req.Properties),
					"properties":   req.Properties,
				}
				if s.JSONOut {
					return command.EmitJSON(output.OK(s.TraceID, "records create dry-run OK", payload))
				}
				fmt.Printf("dry-run create collection=%s db_id=%s properties=%d\n", req.Collection, req.DBID, len(req.Properties))
				return nil
			}
			res, err := s.Runtime.Records.Create(cmd.Context(), req)
			if err != nil {
				return err
			}
			if s.JSONOut {
				return command.EmitJSON(output.OK(s.TraceID, "record created", res))
			}
			fmt.Printf("created id=%s\n", res.ID)
			return nil
		},
		Flags: []command.FlagSpec{
			{Type: command.FlagString, Name: "data", Usage: "Record properties JSON object", Target: &rawData},
			{Type: command.FlagString, Name: "data-file", Usage: "Path to JSON file containing record properties object", Target: &rawFile},
			{Type: command.FlagBool, Name: "validate-schema", Usage: "Validate property keys against Notion database schema", Target: &validateSchema, Default: true},
			{Type: command.FlagBool, Name: "dry-run", Usage: "Validate and preview create without mutation", Target: &dryRun, Default: false},
		},
	}
}

func buildUpdateSpec(s *command.State) command.CommandSpec {
	var req recordsdomain.UpdateRequest
	var rawData, rawFile string
	var validateSchema, dryRun bool
	var ifLastEdited string

	return command.CommandSpec{
		Use:   "update <page-id>",
		Short: "Update one Notion page record properties",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, err := parseProperties(rawData, rawFile)
			if err != nil {
				return err
			}
			req.ID = args[0]
			req.Properties = props
			if ifLastEdited != "" {
				if err := ensureLastEditedMatch(cmd, s, req.ID, ifLastEdited); err != nil {
					return err
				}
			}
			if validateSchema {
				dbID := resolveCollectionDBID(s, req.Collection)
				if strings.TrimSpace(req.Collection) == "" || strings.TrimSpace(dbID) == "" {
					return perr.New(perr.KindValidation, "--collection is required when --validate-schema=true for update")
				}
				if err := validatePropertiesAgainstSchema(cmd, s, req.Collection, dbID, props); err != nil {
					return err
				}
			}
			if err := s.Runtime.AgentAuthorizer.Enforce(policy.Request{
				Actor:      s.Actor,
				Command:    "records.update",
				Collection: normalizeCollection(req.Collection),
				Properties: propertyKeys(req.Properties),
			}); err != nil {
				return err
			}
			if dryRun {
				current, err := s.Runtime.Records.Get(cmd.Context(), recordsdomain.GetRequest{ID: req.ID})
				if err != nil {
					return err
				}
				diff := diffProperties(current.Properties, req.Properties)
				payload := map[string]any{
					"action":       "update",
					"id":           req.ID,
					"collection":   req.Collection,
					"propertyKeys": propertyKeys(req.Properties),
					"changes":      diff,
				}
				if s.JSONOut {
					return command.EmitJSON(output.OK(s.TraceID, "records update dry-run OK", payload))
				}
				fmt.Printf("dry-run update id=%s changes=%d\n", req.ID, len(diff))
				return nil
			}
			res, err := s.Runtime.Records.Update(cmd.Context(), req)
			if err != nil {
				return err
			}
			if s.JSONOut {
				return command.EmitJSON(output.OK(s.TraceID, "record updated", res))
			}
			fmt.Printf("updated id=%s\n", res.ID)
			return nil
		},
		Flags: []command.FlagSpec{
			{Type: command.FlagString, Name: "collection", Usage: "Collection alias or db id (required for agent policy writes)", Target: &req.Collection},
			{Type: command.FlagString, Name: "data", Usage: "Record properties JSON object", Target: &rawData},
			{Type: command.FlagString, Name: "data-file", Usage: "Path to JSON file containing record properties object", Target: &rawFile},
			{Type: command.FlagBool, Name: "validate-schema", Usage: "Validate property keys against Notion database schema", Target: &validateSchema, Default: false},
			{Type: command.FlagBool, Name: "dry-run", Usage: "Validate and preview update without mutation", Target: &dryRun, Default: false},
			{Type: command.FlagString, Name: "if-last-edited", Usage: "Require exact RFC3339 last_edited_time before updating", Target: &ifLastEdited},
		},
	}
}

func buildArchiveSpec(s *command.State) command.CommandSpec {
	var req recordsdomain.ArchiveRequest
	var confirm, dryRun bool
	var ifLastEdited string

	return command.CommandSpec{
		Use:   "archive <page-id>",
		Short: "Archive one Notion page record",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return perr.New(perr.KindValidation, "--confirm is required to archive a record")
			}
			req.ID = args[0]
			if ifLastEdited != "" {
				if err := ensureLastEditedMatch(cmd, s, req.ID, ifLastEdited); err != nil {
					return err
				}
			}
			if err := s.Runtime.AgentAuthorizer.Enforce(policy.Request{
				Actor:      s.Actor,
				Command:    "records.archive",
				Collection: normalizeCollection(req.Collection),
			}); err != nil {
				return err
			}
			if dryRun {
				payload := map[string]any{
					"action":     "archive",
					"id":         req.ID,
					"collection": req.Collection,
				}
				if s.JSONOut {
					return command.EmitJSON(output.OK(s.TraceID, "records archive dry-run OK", payload))
				}
				fmt.Printf("dry-run archive id=%s\n", req.ID)
				return nil
			}
			if err := s.Runtime.Records.Archive(cmd.Context(), req); err != nil {
				return err
			}
			if s.JSONOut {
				return command.EmitJSON(output.OK(s.TraceID, "record archived", map[string]any{"id": req.ID}))
			}
			fmt.Printf("archived id=%s\n", req.ID)
			return nil
		},
		Flags: []command.FlagSpec{
			{Type: command.FlagString, Name: "collection", Usage: "Collection alias or db id (required for agent policy writes)", Target: &req.Collection},
			{Type: command.FlagBool, Name: "confirm", Usage: "Confirm archive mutation", Target: &confirm, Default: false},
			{Type: command.FlagBool, Name: "dry-run", Usage: "Validate and preview archive without mutation", Target: &dryRun, Default: false},
			{Type: command.FlagString, Name: "if-last-edited", Usage: "Require exact RFC3339 last_edited_time before archiving", Target: &ifLastEdited},
		},
	}
}
