# Invoice Module

Path: `internal/invoice`

## Purpose
Domain orchestration for invoice generation, listing, preview, validation, and idempotency.

## Key Flows
- `Create`: load draft, validate, compute totals, render, store artifact, persist idempotency
- `List`: query source and map summary payloads
- `Preview`: render HTML/PDF without invoice status mutation

## Contracts
- Input/Output DTOs: `model.go`
- Public interface: `api.go`
- Workflows: `workflow_*.go`
- Validation rules: `validate.go`

## Idempotency
- Key: `sha256(page_id:last_edited_time:template_version)`
- Store: local JSON file (`invoice.idempotency_store`)
- File permissions are restricted (`0600`) with parent dir `0700`.

## Source Strategy
- Primary: Notion (`invoice.source=notion`)
- Optional fallback: local JSON fixture (`invoice.fallback_file`)

## Safety Controls
- Notion mutations are gated by config:
  - `invoice.allow_notion_mutations`
  - `invoice.require_mutation_confirm`
- Create command requires `--confirm` when `--upload-notion` is used and confirmation is enabled.
- Renderer hardening defaults:
  - JavaScript disabled
  - Chrome sandbox enabled by default (`renderer_no_sandbox: false`)
  - Rendered HTML size cap (`invoice.max_render_html_bytes`)
- Optional output write boundary:
  - `invoice.output_base_dir` enforces artifact path containment.
- Agent command policy hooks:
  - command allowlist, invoice ID regex, and list max limit are enforced for `--actor agent`.
