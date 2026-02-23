# Records Module

Path: `internal/modules/records`, `internal/records`

## Purpose
Provide generic access to Notion records outside invoice-specific workflows.

## Commands
- `biz records list <collection-or-db-id> [--limit N] [--cursor C]`
- `biz records get <page-id>`
- `biz records schema <collection-or-db-id>`
- `biz records create <collection-or-db-id> --data '<json properties>' | --data-file <path> [--dry-run]`
- `biz records update <page-id> --collection <collection-or-db-id> --data '<json properties>' | --data-file <path> [--dry-run] [--if-last-edited RFC3339]`
- `biz records archive <page-id> --collection <collection-or-db-id> --confirm [--dry-run] [--if-last-edited RFC3339]`

## Collection Resolution
- `invoices` or `invoice` -> uses `notion.invoice_db_id` from config
- Any key in `notion.collections` -> maps to configured database id
- Any other value is treated as a raw Notion database id for `list`

## Output
- `--json` returns the standard envelope with raw Notion page `properties`
- Plain output prints record ids and pagination cursor

## Write Safety
- `records create` validates property names against Notion schema by default (`--validate-schema=true`).
- `records update` can validate schema when `--collection` is provided (`--validate-schema=true`).
- `--dry-run` previews mutations without writing.
- `--if-last-edited` prevents updates/archives if the page changed since your expected timestamp.

## Agent Safety
For `--actor agent`, writes are controlled by `agent_policy`:

- `allowed_commands` must include `records.create`, `records.update`, `records.archive`
- `records_allowed_collections` must include the provided collection
- `records_allowed_properties` must include all top-level property names in `--data`

## Enablement
The module is loaded only when configured in:

```yaml
modules:
  enabled: [invoice, records]
```
