# Tasks: Issue #2 Config Core

## Phase 1: Guard Alignment (Batch 1)

- [x] 1.1 Update `bootstrap_guard_test.go` to remove `config` from doc-only package enforcement, while keeping guards for bootstrap-only packages (including `config/viper`).
- [x] 1.2 Add/adjust guard assertions in `bootstrap_guard_test.go` so `config/` implementation files are allowed for this change, without broadening unrelated package allowances.

## Phase 2: Core Types and Constructors (Batch 2)

- [x] 2.1 Create `config/types.go` with typed model: `Config`, `ServerConfig`, `SecurityConfig`, `AuthConfig`, `JWTConfig`, `CookieConfig`, `LoginConfig`, `OAuth2Config`, and generic `OAuth2ProviderConfig`.
- [x] 2.2 Implement `config/constructors.go` with `NewConfig` and focused `New*Config` helpers using dependency injection inputs (no globals) and useful zero/default behavior.
- [x] 2.3 Ensure constructor outputs are deterministic and side-effect free (no env/filesystem/time reads).

## Phase 3: Errors, Validation, and Normalization (Batch 3)

- [x] 3.1 Create `config/errors.go` with stable sentinels (`ErrInvalidConfig`, `ErrRequired`, `ErrOutOfRange`, `ErrInvalidEmail`) and `ValidationError{Path, Err}` with `Unwrap()`.
- [x] 3.2 Create `config/validate.go` with `ValidateConfig(cfg Config) (Config, error)` orchestrating focused validators: server, jwt, cookie, login, oauth2.
- [x] 3.3 Implement login email normalization in validation flow (trim spaces + lowercase) before success is returned.
- [x] 3.4 Wrap validation failures with contextual path info and preserve `errors.Is`/`errors.As` compatibility.

## Phase 4: Unit Tests (Batch 4)

- [x] 4.1 Add `config/constructors_test.go` table-driven tests for deterministic constructor behavior, defaults, and zero-value usefulness.
- [x] 4.2 Add `config/validate_test.go` valid-path tests covering baseline compliant config for server/JWT/cookie/login/oauth2.
- [x] 4.3 Add `config/validate_test.go` invalid-path tests for missing/invalid values asserting sentinels via `errors.Is` and `ValidationError` via `errors.As`.
- [x] 4.4 Add normalization tests in `config/validate_test.go` verifying normalized login emails are trimmed + lowercased in returned config.
- [x] 4.5 Add `config/errors_test.go` to verify `ValidationError` path metadata and wrapping behavior.

## Phase 5: Docs, Snippets, and Verification (Batch 5)

- [x] 5.1 Update `config/doc.go` with package contract: typed core, panic-free validation, no globals, no adapter coupling.
- [x] 5.2 Update `README.md` with minimal snippet showing construction + `ValidateConfig` usage and error assertions.
- [x] 5.3 Verify acceptance criteria by running `go test ./...` and checking spec alignment: typed model present, guard first-pass, normalization works, assertable errors, and no `viper` usage in core files.
- [x] 5.4 Record verification notes in change context (map each spec scenario in `openspec/changes/issue-2-config-core/specs/**/spec.md` to test coverage).

## Verification Notes

- `config-core` / **Model supports issue baseline domains**: covered by compile-time usage in `config/types.go` and construction in `config/constructors_test.go` (`TestNewConfigIsDeterministicAndCopiesDependencies`) and `config/validate_test.go` fixture.
- `config-core` / **Provider model remains generic**: covered by `OAuth2ProviderConfig` and `OAuth2Config` in `config/types.go` and exercised in `config/validate_test.go` provider cases.
- `config-core` / **Valid inputs construct config deterministically**: covered by `config/constructors_test.go` (`TestNewServerConfig`, `TestNewJWTConfig`, `TestNewConfigIsDeterministicAndCopiesDependencies`, `TestNewOAuth2ProviderConfigCopiesScopes`).
- `config-core` / **Invalid inputs do not panic**: covered by `ValidateConfig` error-return behavior in `config/validate_test.go` invalid table tests and `config/errors_test.go` wrapping checks.
- `config-core` / **Baseline validation succeeds for compliant config**: covered by `TestValidateConfigValid` in `config/validate_test.go`.
- `config-core` / **Baseline validation reports assertable failures**: covered by `TestValidateConfigInvalid` + `errors.Is`/`errors.As` assertions in `config/validate_test.go` and `config/errors_test.go`.
- `config-core` / **Login normalization trim+lowercase**: covered by `TestValidateConfigNormalizesLoginEmail` in `config/validate_test.go`.
- `config-core` / **Core package runs without adapter dependencies**: core files `config/*.go` import only standard library; adapter remains isolated under `config/viper`.
- `config-core` / **Repeated executions are side-effect free**: constructor copy/determinism and pure validation path covered in constructor/validation tests.
- `framework-bootstrap` / **Bootstrap guard allows approved config evolution**: covered by updated `bootstrap_guard_test.go` guard list and `TestConfigPackageCanEvolveBeyondBootstrapDocs`.
