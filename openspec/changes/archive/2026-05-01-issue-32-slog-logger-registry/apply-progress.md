# Apply Progress: issue-32-slog-logger-registry

## Change
- **Name**: `issue-32-slog-logger-registry`
- **Mode**: Standard (Strict TDD disabled by project config)
- **Artifact store**: hybrid

## Completed Tasks

### Batch A (Phase 1)
- [x] 1.1 Add config tests for logging defaults/shape
- [x] 1.2 Add logging config model and constructor defaults
- [x] 1.3 Add validation tests for logging format/levels/logger keys
- [x] 1.4 Implement logging validation and normalization
- [x] 1.5 Create core logging contracts and precedence resolver

### Batch B (Phase 2)
- [x] 2.1 Add slog registry stability/isolation tests
- [x] 2.2 Implement slog registry with mutex cache and deterministic settings resolution
- [x] 2.3 Add json/text output + warn filtering tests
- [x] 2.4 Implement slog logger adapter and noop logger

### Batch C (Phase 3)
- [x] 3.1 Add app lifecycle tests for `GetLogger(name)`
- [x] 3.2 Implement `Application.GetLogger(name)` and runtime wiring
- [x] 3.3 Add per-application isolation and concurrent `GetLogger` tests
- [x] 3.4 Update app docs for logger lifecycle and constraints

### Batch D (Phase 4)
- [x] 4.1 Add viper mapping/loader tests for logging keys and env precedence
- [x] 4.2 Implement logging mapping + env overrides (`logging.loggers.<name>.*`)
- [x] 4.3 Add wrapped validation failure-path coverage for invalid logging values

### Batch E (Phase 5)
- [x] 5.1 Update README logging contract/config/precedence/examples
- [x] 5.2 Update docs home/migration/release checklist with collector-first Loki guidance
- [x] 5.3 Extend docs-contract assertions for logging API/config/output sync
- [x] 5.4 Run verification (`go test ./...`, `go build ./...`) and add migration note summary

## Verification Evidence
- `go test ./...` ✅
- `go build ./...` ✅

## Design Alignment
- Implemented core/adapter split via `logging` + `logging/slog` packages.
- Preserved deterministic precedence semantics (root + per-logger overrides).
- Enforced per-application runtime isolation and same-name logger stability.

## Deviations
- None — implementation matches design and spec requirements.

## Remaining
- None. All tasks in `tasks.md` are marked complete.
