# Tasks: issue-28-http-server-timeout-header-size

## Phase 1: Config Core Foundation

- [x] 1.1 Modify `config/types.go` to extend `ServerConfig` with `ReadTimeout`, `WriteTimeout`, and `MaxHeaderBytes`.
- [x] 1.2 Modify `config/constructors.go` so `NewServerConfig` sets defaults (`10s`, `10s`, `1048576`) when runtime-limit inputs are omitted.
- [x] 1.3 Update `config/constructors.go` call sites impacted by constructor signature changes to keep compilation green.
- [x] 1.4 Extend `config/validate.go` (`validateServer`) to reject zero/negative timeout and header-size values with contextual validation errors.

## Phase 2: Adapter and Runtime Wiring

- [x] 2.1 Modify `config/viper/mapping.go` to map `server.read-timeout`, `server.write-timeout`, and `server.max-header-bytes` into core config.
- [x] 2.2 Modify `config/viper/loader.go` to support env overrides `COMMON_FWK_SERVER_READ_TIMEOUT`, `COMMON_FWK_SERVER_WRITE_TIMEOUT`, and `COMMON_FWK_SERVER_MAX_HEADER_BYTES`.
- [x] 2.3 In `config/viper/loader.go`, parse env runtime limits explicitly (`time.ParseDuration`, `strconv.Atoi`) and return adapter-typed decode/load errors on invalid values.
- [x] 2.4 Modify `app/application.go` so `UseServer()` applies `a.cfg.Server.ReadTimeout`, `WriteTimeout`, and `MaxHeaderBytes` to the underlying `http.Server` while preserving fluent chaining.

## Phase 3: Automated Tests (Table-Driven, Deterministic)

- [x] 3.1 Extend `config/constructors_test.go` with table-driven `t.Run` cases for default application and explicit runtime-limit preservation.
- [x] 3.2 Extend `config/validate_test.go` with table-driven positive/zero/negative runtime-limit cases and assertable validation error checks.
- [x] 3.3 Extend `config/viper/mapping_test.go` to verify deterministic file-key mapping for all three runtime-limit fields.
- [x] 3.4 Extend `config/viper/loader_test.go` to verify env precedence, deterministic override behavior, and typed failures for invalid duration/int inputs.
- [x] 3.5 Extend `app/application_test.go` to verify fluent chain behavior and `UseServer()` runtime-limit wiring (explicit values and defaults).

## Phase 4: Documentation Synchronization

- [x] 4.1 Update `README.md` with server runtime-limit keys, defaults, env variable names, and at least one configuration example.
- [x] 4.2 Update `docs/home.md` and related config docs under `docs/` to include the same keys/defaults/example and keep wording aligned with README.

## Phase 5: Verification and Readiness

- [x] 5.1 Run `go test ./...` and confirm all updated tests pass.
- [x] 5.2 Run `go build ./...` to verify no integration/build regressions.
- [x] 5.3 Verify completion against spec scenarios (defaults, explicit preservation, validation failures, adapter mapping/overrides, typed decode failures, app wiring, docs sync).
