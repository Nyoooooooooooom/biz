# Platform Modules

Path: `internal/platform`

## Packages
- `config`: viper-backed config loader
- `errors`: typed domain/platform errors + exit code mapping
- `output`: stable API envelope (`v1`)
- `log`: zap logger factory
- `clock`: testable time source
- `id`: trace id generation
- `audit` (outside `internal/platform`): signed hash-chained JSONL audit writer

## Rule of Thumb
Platform packages must remain domain-agnostic and reusable.
