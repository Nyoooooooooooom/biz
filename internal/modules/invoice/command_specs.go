package invoice

import (
	"fmt"

	"biz/internal/command"
	"biz/internal/invoice"
	"biz/internal/platform/output"
	"biz/internal/policy"
	"github.com/spf13/cobra"
)

func buildInvoiceCommand(s *command.State) *cobra.Command {
	listSpec := buildListSpec(s)
	createSpec := buildCreateSpec(s)
	previewSpec := buildPreviewSpec(s)
	return command.BuildCommand(command.CommandSpec{
		Use:     "invoice",
		Aliases: []string{"inv"},
		Short:   "Invoice operations",
		RunE:    listSpec.RunE,
		Commands: []command.CommandSpec{
			createSpec,
			listSpec,
			previewSpec,
		},
	})
}

func buildListSpec(s *command.State) command.CommandSpec {
	var req invoice.ListRequest
	return command.CommandSpec{
		Use:     "list [status]",
		Aliases: []string{"ls"},
		Short:   "List invoices",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				req.Status = args[0]
			}
			if req.Status == "" {
				req.Status = s.Runtime.Config.Invoice.DefaultStatusQuery
			}
			if req.Limit <= 0 {
				req.Limit = s.Runtime.Config.Invoice.DefaultListLimit
			}
			if err := s.Runtime.AgentAuthorizer.Enforce(policy.Request{
				Actor:     s.Actor,
				Command:   "invoice.list",
				ListLimit: req.Limit,
			}); err != nil {
				return err
			}
			req.TraceID = s.TraceID
			res, err := s.Runtime.Invoice.List(cmd.Context(), req)
			if err != nil {
				return err
			}
			if s.JSONOut {
				return command.EmitJSON(output.OK(s.TraceID, "invoices listed", res))
			}
			for _, item := range res.Items {
				fmt.Printf("%s %s %s %.2f %s\n", item.InvoiceNumber, item.ClientName, item.Status, item.Total, item.Currency)
			}
			if res.NextCursor != "" {
				fmt.Printf("next_cursor=%s\n", res.NextCursor)
			}
			return nil
		},
		Flags: []command.FlagSpec{
			{Type: command.FlagString, Name: "status", Usage: "Filter by status", Target: &req.Status},
			{Type: command.FlagInt, Name: "limit", Usage: "Maximum results", Target: &req.Limit, Default: 0},
			{Type: command.FlagString, Name: "cursor", Usage: "Pagination cursor", Target: &req.Cursor},
		},
	}
}

func buildCreateSpec(s *command.State) command.CommandSpec {
	var req invoice.CreateRequest
	return command.CommandSpec{
		Use:     "create <invoice_id>",
		Aliases: []string{"cr", "new"},
		Short:   "Create invoice PDF",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req.ID = args[0]
			if req.OutDir == "" {
				req.OutDir = s.Runtime.Config.Invoice.OutputDir
			}
			if req.Source == "" {
				req.Source = s.Runtime.Config.Invoice.Source
			}
			if !cmd.Flags().Changed("upload-notion") {
				req.UploadNotion = s.Runtime.Config.Invoice.DefaultUploadNotion
			}
			if err := s.Runtime.AgentAuthorizer.Enforce(policy.Request{
				Actor:     s.Actor,
				Command:   "invoice.create",
				InvoiceID: req.ID,
			}); err != nil {
				return err
			}
			req.TraceID = s.TraceID
			res, err := s.Runtime.Invoice.Create(cmd.Context(), req)
			if err != nil {
				return err
			}
			if s.JSONOut {
				return command.EmitJSON(output.OK(s.TraceID, "invoice created", res))
			}
			fmt.Printf("created invoice %s -> %s\n", res.InvoiceNumber, res.PDFPath)
			return nil
		},
		Flags: []command.FlagSpec{
			{Type: command.FlagString, Name: "out", Usage: "Output directory", Target: &req.OutDir},
			{Type: command.FlagBool, Name: "upload-notion", Usage: "Upload PDF back to Notion", Target: &req.UploadNotion, Default: false},
			{Type: command.FlagBool, Name: "confirm", Usage: "Confirm mutation when using --upload-notion", Target: &req.Confirm, Default: false},
			{Type: command.FlagString, Name: "source", Usage: "Data source (local|notion)", Target: &req.Source},
			{Type: command.FlagString, Name: "source-file", Usage: "Local JSON source file", Target: &req.SourceFile},
		},
	}
}

func buildPreviewSpec(s *command.State) command.CommandSpec {
	var req invoice.PreviewRequest
	return command.CommandSpec{
		Use:     "preview <invoice_id>",
		Aliases: []string{"pv", "show"},
		Short:   "Preview invoice as HTML/PDF",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req.ID = args[0]
			if req.Format == "" {
				req.Format = s.Runtime.Config.Invoice.DefaultPreviewFormat
			}
			if err := s.Runtime.AgentAuthorizer.Enforce(policy.Request{
				Actor:     s.Actor,
				Command:   "invoice.preview",
				InvoiceID: req.ID,
			}); err != nil {
				return err
			}
			req.TraceID = s.TraceID
			res, err := s.Runtime.Invoice.Preview(cmd.Context(), req)
			if err != nil {
				return err
			}
			if s.JSONOut {
				return command.EmitJSON(output.OK(s.TraceID, "invoice preview generated", res))
			}
			fmt.Printf("preview generated: %s (%s, %d bytes)\n", res.Path, res.MimeType, res.SizeBytes)
			return nil
		},
		Flags: []command.FlagSpec{
			{Type: command.FlagString, Name: "format", Usage: "Preview format (html|pdf)", Target: &req.Format, Default: ""},
		},
	}
}
