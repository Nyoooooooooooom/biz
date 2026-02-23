# Contributing

## Development
```bash
go test ./...
go vet ./...
```

## Code Style
- Keep command handlers thin.
- Place business logic in `internal/<domain>`.
- Use typed errors from `internal/platform/errors`.
- Preserve machine output contract in `internal/platform/output`.

## Pull Request Checklist
- [ ] tests added/updated
- [ ] docs updated (`README` or `docs/modules/*`)
- [ ] no secrets in committed files
- [ ] no local artifacts committed (`config.yaml`, `invoices/`, caches)
