# Release Guide

## Preconditions

- `go test ./...` passes.
- `go vet ./...` passes.
- `docker build -t biz:release .` passes.
- `config.yaml` is excluded from commits.

## Release Steps

1. Update `CHANGELOG.md` under a new version heading.
2. Tag the release:
   - `git tag -a v0.1.0 -m "v0.1.0"`
3. Push commits and tags:
   - `git push origin main --tags`
4. Create GitHub release notes from `CHANGELOG.md`.

## Post-Release Smoke

- `docker run --rm biz:release --help`
- `docker run --rm -v $(pwd)/config.yaml:/app/config.yaml biz:release --config /app/config.yaml invoice list --json`
