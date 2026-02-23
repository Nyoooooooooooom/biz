# CLI Module

Path: `internal/command`, `internal/modules/*`

## Purpose
Expose stable command surfaces for humans and agents through composable modules.

## Commands
- `biz invoice list`
- `biz invoice preview`
- `biz invoice create`
- `biz doctor`

## Principles
- Always support `--json` envelope output.
- Write machine-readable output to `stdout`.
- Write diagnostic/errors to `stderr`.
- Keep module handlers thin; delegate business logic to domain services.

## Extension Pattern
Add a new module by:
1. creating `internal/modules/<name>/module.go` implementing `command.Module`,
2. defining commands using `command.CommandSpec` + `command.FlagSpec`,
3. delegating to fragment/domain services,
4. registering the module in `cmd/biz/main.go`.
