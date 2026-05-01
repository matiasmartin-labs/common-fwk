# Apply Progress: issue-28-http-server-timeout-header-size

## Mode

Standard (strict_tdd: false)

## Completed Tasks (Merged Cumulative State)

- [x] 1.1 Modify `config/types.go` to extend `ServerConfig` with `ReadTimeout`, `WriteTimeout`, and `MaxHeaderBytes`.
- [x] 1.2 Modify `config/constructors.go` so `NewServerConfig` sets defaults (`10s`, `10s`, `1048576`) when runtime-limit inputs are omitted.
- [x] 1.3 Update `config/constructors.go` call sites impacted by constructor signature changes to keep compilation green.
- [x] 1.4 Extend `config/validate.go` (`validateServer`) to reject zero/negative timeout and header-size values with contextual validation errors.
- [x] 2.1 Modify `config/viper/mapping.go` to map `server.read-timeout`, `server.write-timeout`, and `server.max-header-bytes` into core config.
- [x] 2.2 Modify `config/viper/loader.go` to support env overrides `COMMON_FWK_SERVER_READ_TIMEOUT`, `COMMON_FWK_SERVER_WRITE_TIMEOUT`, and `COMMON_FWK_SERVER_MAX_HEADER_BYTES`.
- [x] 2.3 In `config/viper/loader.go`, parse env runtime limits explicitly (`time.ParseDuration`, `strconv.Atoi`) and return adapter-typed decode/load errors on invalid values.
- [x] 2.4 Modify `app/application.go` so `UseServer()` applies `a.cfg.Server.ReadTimeout`, `WriteTimeout`, and `MaxHeaderBytes` to the underlying `http.Server` while preserving fluent chaining.
- [x] 3.1 Extend `config/constructors_test.go` with table-driven `t.Run` cases for default application and explicit runtime-limit preservation.
- [x] 3.2 Extend `config/validate_test.go` with table-driven positive/zero/negative runtime-limit cases and assertable validation error checks.
- [x] 3.3 Extend `config/viper/mapping_test.go` to verify deterministic file-key mapping for all three runtime-limit fields.
- [x] 3.4 Extend `config/viper/loader_test.go` to verify env precedence, deterministic override behavior, and typed failures for invalid duration/int inputs.
- [x] 3.5 Extend `app/application_test.go` to verify fluent chain behavior and `UseServer()` runtime-limit wiring (explicit values and defaults).
- [x] 4.1 Update `README.md` with server runtime-limit keys, defaults, env variable names, and at least one configuration example.
- [x] 4.2 Update `docs/home.md` and related config docs under `docs/` to include the same keys/defaults/example and keep wording aligned with README.
- [x] 5.1 Run `go test ./...` and confirm all updated tests pass.
- [x] 5.2 Run `go build ./...` to verify no integration/build regressions.
- [x] 5.3 Verify completion against spec scenarios (defaults, explicit preservation, validation failures, adapter mapping/overrides, typed decode failures, app wiring, docs sync).

## Continuation Batch — Verify Findings Remediation

- [x] Added automated docs contract test evidence for required runtime-limit contract in both `README.md` and `docs/home.md`.
- [x] Strengthened env override determinism assertions for runtime-limit override path.
- [x] Re-ran `go test ./...` and `go build ./...` after remediation.

## Files Changed In This Continuation Batch

- `bootstrap_guard_test.go`
  - Added `TestDocsRuntimeLimitsContract` with table-driven `t.Run` cases for `README.md` and `docs/home.md`.
  - Asserts required runtime-limit keys (`read-timeout`, `write-timeout`, `max-header-bytes`), defaults (`10s`, `1048576`), env vars (`COMMON_FWK_SERVER_*`), and YAML example markers are present.
- `config/viper/loader_test.go`
  - Updated `TestLoadEnvOverrideSemantics` to load with `EnvOverride=true` twice and assert deterministic output using `reflect.DeepEqual` for identical env snapshots.

## Validation Evidence (Current Batch)

### `go test ./...`

Passed:

- `github.com/matiasmartin-labs/common-fwk`
- `github.com/matiasmartin-labs/common-fwk/app`
- `github.com/matiasmartin-labs/common-fwk/config`
- `github.com/matiasmartin-labs/common-fwk/config/viper`
- `github.com/matiasmartin-labs/common-fwk/errors`
- `github.com/matiasmartin-labs/common-fwk/http/gin`
- `github.com/matiasmartin-labs/common-fwk/security/claims`
- `github.com/matiasmartin-labs/common-fwk/security/jwt`
- `github.com/matiasmartin-labs/common-fwk/security/keys`

No tests:

- `github.com/matiasmartin-labs/common-fwk/security`

### `go build ./...`

Passed (no build errors).

## Spec Scenario Traceability (Updated)

- Config core defaults + explicit preservation: covered via constructor changes and `config/constructors_test.go` table-driven cases.
- Runtime-limit validation failures: covered in `config/validate.go` and `config/validate_test.go` zero/negative cases.
- Viper mapping and env precedence: covered in `config/viper/mapping.go`, `config/viper/loader.go`, and expanded tests.
- Typed adapter failures for env decoding: covered by explicit parse branches and `TestLoadEnvOverrideTypedFailuresForServerRuntimeLimits`.
- App runtime wiring + fluent chaining: covered in `app/application.go` and `app/application_test.go`.
- Docs runtime-limit contract scenario: now backed by automated test evidence in `bootstrap_guard_test.go` (`TestDocsRuntimeLimitsContract`).
- Env override precedence + deterministic behavior scenario: now explicitly asserts deterministic repeated `EnvOverride=true` output in `config/viper/loader_test.go` (`TestLoadEnvOverrideSemantics`).

## Status

18/18 original tasks remain complete. Verify-blocking findings addressed. Ready for re-verify.
