# Tasks: Issue #3 Optional Config Viper Adapter

## Phase 1: Foundation (adapter contracts and deterministic options)

- [x] 1.1 Create `config/viper/options.go` with `Options` defaults and explicit semantics for `ConfigPath`, `ConfigType`, `EnvPrefix`, `EnvOverride`, `ExpandEnv`; validate with `go test ./config/viper -run TestOptions`. [Req: Deterministic option semantics]
- [x] 1.2 Create `config/viper/errors.go` with `LoadError`, `DecodeError`, `MappingError`, `ValidationError` and `Unwrap()` support; validate with `go test ./config/viper -run TestErrors`. [Req: Explicit mapping and typed adapter errors]
- [x] 1.3 Create `config/viper/mapping.go` adapter-local raw structs and explicit raw→`config.Config` mapping helpers that return `MappingError` on invalid raw values; validate with `go test ./config/viper -run TestMapping`. [Req: Explicit mapping and typed adapter errors]

## Phase 2: Core implementation (loader orchestration)

- [x] 2.1 Create `config/viper/loader.go` implementing panic-free `Load(opts Options) (config.Config, error)` using a fresh Viper instance per call, deterministic option application, and contextual stage errors; validate with `go test ./config/viper -run TestLoad`. [Req: Loader API contract]
- [x] 2.2 In `config/viper/loader.go`, call `config.ValidateConfig` after mapping and wrap failures with `ValidationError` while preserving `errors.Is/As` for core sentinels; validate with `go test ./config/viper -run TestLoadWrapsCoreValidation`. [Req: Mandatory post-load core validation, Delta: config-core]
- [x] 2.3 Finalize and codify `ConfigType` behavior (explicit vs extension inference) in `options.go`/`loader.go` and enforce it with deterministic error paths; validate with focused table tests. [Req: Deterministic option semantics]

## Phase 3: Testing and verification

- [x] 3.1 Create `config/viper/loader_test.go` success-path tests using temp config files and isolated env snapshots; assert repeated identical inputs return identical `config.Config`. [Req: Loader API contract, Deterministic option semantics]
- [x] 3.2 In `config/viper/loader_test.go`, add failure tests for unreadable/missing file and malformed config content with stage-typed errors (`LoadError`, `DecodeError`). [Req: Loader API contract, typed adapter errors]
- [x] 3.3 In `config/viper/loader_test.go`, add env precedence and expansion tests (`EnvOverride` on/off, `ExpandEnv` on/off) and deterministic expansion assertions. [Req: Deterministic option semantics, Environment expansion determinism]
- [x] 3.4 Create `config/viper/mapping_test.go` mapping-stage failure cases and stable provider/scopes mapping assertions. [Req: Explicit mapping and typed adapter errors]
- [x] 3.5 Create `config/viper/errors_test.go` asserting `errors.As` for adapter stages and `errors.Is`/`errors.As` for wrapped core validation classes (e.g., `config.ErrInvalidConfig`). [Req: Mandatory post-load core validation, Delta: config-core]

## Phase 4: Integration and documentation

- [x] 4.1 Update `go.mod` and `go.sum` to add `github.com/spf13/viper` for adapter package usage only; validate with `go test ./...` and `go build ./...`.
- [x] 4.2 Update `config/viper/doc.go` and `README.md` with adapter contract, precedence rules, deterministic behavior, and typed error usage example; validate by matching documented behavior to test coverage.
