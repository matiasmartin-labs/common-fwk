# Tasks: issue-32-slog-logger-registry

AC map: A1 `GetLogger` deterministic; A2 precedence/filtering; A3 fields+formats; A4 isolation; A5 docs+Loki guidance.

## Phase 1: Foundation — contracts and config model

- [x] 1.1 RED: Add failing config tests for logging defaults/shape in `config/constructors_test.go` and `config/types.go` expectations. (AC: A2, A3 | Test: root defaults and per-logger override typing)
- [x] 1.2 GREEN: Add `LoggingConfig` and per-logger override structs in `config/types.go`; wire defaults in `config/constructors.go` (`enabled=true`,`level=info`,`format=json`). (AC: A2, A3)
- [x] 1.3 RED: Add failing validation tests in `config/validate_test.go` for invalid `logging.format`, invalid levels, and logger-key constraints. (AC: A2 | Test: contextual assertable errors)
- [x] 1.4 GREEN: Implement logging validation in `config/validate.go` for `json|text`, level enum, and override checks. (AC: A2)
- [x] 1.5 Create core contracts in `logging/logger.go` and `logging/registry.go` (interfaces, level helpers, precedence resolver inputs). (AC: A1, A2)

## Phase 2: Slog adapter registry implementation

- [x] 2.1 RED: Add failing unit tests in `logging/slog/registry_test.go` for same-name stability and cross-name isolation. (AC: A1, A4)
- [x] 2.2 GREEN: Implement `logging/slog/registry.go` with `sync.RWMutex` + per-name cache and deterministic effective setting resolution. (AC: A1, A2, A4)
- [x] 2.3 RED: Add failing output/filter tests in `logging/slog/logger_test.go` for `json|text` fields (`logger`,`ts`,`level`,`msg`) and warn-level filtering. (AC: A2, A3)
- [x] 2.4 GREEN: Implement `logging/slog/logger.go` facade and `logging/slog/noop.go` disabled logger; satisfy filtering and field contract. (AC: A2, A3)

## Phase 3: App lifecycle and integration wiring

- [x] 3.1 RED: Extend `app/application_test.go` with failing scenarios: pre-bootstrap error, empty-name error, post-bootstrap stable logger instance. (AC: A1)
- [x] 3.2 GREEN: Update `app/application.go` to own registry runtime and expose `GetLogger(name)` with deterministic contextual errors. (AC: A1)
- [x] 3.3 Add concurrency/integration tests in `app/application_test.go` for per-application isolation and parallel `GetLogger` behavior. (AC: A4)
- [x] 3.4 Update `app/doc.go` API docs for `GetLogger(name)` lifecycle/constraints. (AC: A5)

## Phase 4: Viper mapping and precedence behavior

- [x] 4.1 RED: Add failing adapter tests in `config/viper/mapping_test.go` + `loader_test.go` for file mapping of `logging.*` and env override precedence. (AC: A2)
- [x] 4.2 GREEN: Implement `config/viper/mapping.go` and `config/viper/loader.go` logging key mapping (`logging.loggers.<name>.*`) and deterministic env-over-file precedence. (AC: A2)
- [x] 4.3 Add failure-path tests for wrapped core validation errors on invalid mapped logging values. (AC: A2 | Test: assertable error class preserved)

## Phase 5: Docs sync, migration notes, and end-to-end verification

- [x] 5.1 Update `README.md` with logging config keys, precedence matrix, `GetLogger` usage, required output fields, and json/text examples. (AC: A5)
- [x] 5.2 Update `/docs/*`: `docs/home.md`, `docs/migration/auth-provider-ms-v0.1.0.md`, `docs/releases/v0.2.0-checklist.md` with collector-first Loki guidance and structured-field preservation notes. (AC: A5)
- [x] 5.3 Add/extend docs-contract assertions (existing docs verification tests) to enforce logging API/config docs synchronization. (AC: A5)
- [x] 5.4 Run full verification batch: `go test ./...` and `go build ./...`; capture migration note summary in release checklist entry. (AC: A1, A2, A3, A4, A5)

## Apply Batches (recommended continuity)

- Batch A: 1.1–1.5 (config model + contracts)
- Batch B: 2.1–2.4 (slog registry + output behavior)
- Batch C: 3.1–3.4 (app wiring + lifecycle tests)
- Batch D: 4.1–4.3 (viper precedence)
- Batch E: 5.1–5.4 (docs + migration + final verification)
