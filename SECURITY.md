# Security

## Secrets
- Never commit integration tokens to git.
- Keep runtime secrets in environment variables or local-only files.

Recommended env vars:
- `BIZ_NOTION_TOKEN`
- `BIZ_NOTION_INVOICE_DB_ID`

## Sensitive Files
- `config.yaml` is local-only and ignored by git.

## Runtime Hardening Defaults
- Notion mutations are disabled unless explicitly enabled:
  - `invoice.allow_notion_mutations: false`
- Mutating create calls require explicit confirmation:
  - `invoice.require_mutation_confirm: true`
  - CLI: `invoice create ... --upload-notion --confirm`
- Render hardening:
  - `invoice.renderer_disable_javascript: true`
  - `invoice.renderer_no_sandbox: false`
  - `invoice.max_render_html_bytes: 1048576`
- Local artifact hardening:
  - output files are written with `0600` permissions
  - idempotency store is written with `0600` and parent dir `0700`
- optional path boundary via `invoice.output_base_dir`

## Agent Authorization Policy
- Agent constraints are config-driven under `agent_policy` and apply only for `--actor agent`.
- Controls supported:
  - command allowlist (`allowed_commands`)
  - invoice ID regex guard (`invoice_id_regex`)
  - maximum list batch size (`max_list_limit`)

## Signed Audit Logging
- `audit.enabled` activates append-only JSONL events with hash chaining:
  - each event stores `prev_hash` -> `hash`
  - each hash is HMAC-signed with `audit.signing_key`
- Recommended:
  - keep `audit.strict: true`
  - store `audit.signing_key` in env/secret manager, not committed files
  - restrict file perms (`dir_perm: 0700`, `file_perm: 0600`)

## Reporting
If you discover a security issue, open a private report with steps to reproduce and potential impact.
