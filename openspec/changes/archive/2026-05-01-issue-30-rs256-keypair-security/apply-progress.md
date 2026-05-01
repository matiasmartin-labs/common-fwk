# Apply Progress: issue-30-rs256-keypair-security

Mode: Standard (`openspec/config.yaml` sets `strict_tdd: false`; governed checks run via `go test ./...` + `go build ./...`)

## Task Checklist

### Phase 1: Config Foundation
- ✅ 1.1 Extend `config/types.go` with JWT mode-aware fields (`Algorithm`, `RS256Config`) and RS256 source constants while preserving legacy fields.
- ✅ 1.2 Keep `NewJWTConfig(secret, issuer, ttlMinutes)` backward-compatible and add focused RS256 helpers (`NewRS256GeneratedConfig`, `NewRS256PublicPEMConfig`, `NewRS256PrivatePEMConfig`).
- ✅ 1.3 Implement conditional JWT validation for HS256 vs RS256, including unsupported algorithm/source error paths and default algorithm normalization.

### Phase 2: Adapter + Security Core Wiring
- ✅ 2.1 Add RS256 canonical mapping in `config/viper/mapping.go` and compatibility aliases/env override handling in `config/viper/loader.go`.
- ✅ 2.2 Create `security/keys/keypair.go` with in-memory RSA generation and resolver bootstrap for `generated`, `public-pem`, and `private-pem`.
- ✅ 2.3 Update `security/jwt/compat.go` so `FromConfigJWT` branches by algorithm and returns algorithm-constrained options + TTL with contextual errors.
- ✅ 2.4 Add `UseServerSecurityFromConfig()` to `app/application.go` as a thin wrapper: validate config, derive compat options, build validator, and delegate to `UseServerSecurity`.

### Phase 3: Deterministic Tests and Invalid-Mode Coverage
- ✅ 3.1 Expand `config/validate_test.go` with HS256/RS256 valid-invalid matrix coverage.
- ✅ 3.2 Update `config/viper/*_test.go` for RS256 canonical precedence, legacy alias compatibility, env override behavior, and deterministic repeated output.
- ✅ 3.3 Create `security/keys/keypair_test.go` for generation/retrieval success and invalid-classification failures.
- ✅ 3.4 Add `security/jwt/compat_test.go` for HS256/RS256 resolver wiring, method allowlists, and invalid-mode failures.
- ✅ 3.5 Update `app/application_test.go` for `UseServerSecurityFromConfig()` success cases and no-partial-wiring failure behavior.

### Phase 4: Documentation and Verification
- ✅ 4.1 Update `README.md` and `docs/home.md` with HS256 default + RS256 bootstrap semantics and examples.
- ✅ 4.2 Update `docs/migration/auth-provider-ms-v0.1.0.md` with executable HS256→RS256 migration steps and parity checklist.
- ✅ 4.3 Update `docs/releases/v0.2.0-checklist.md` with HS256/RS256 release verification checkpoints.
- ✅ 4.4 Run full verification commands: `go test ./...` and `go build ./...` (green), including doc-contract assertions for release checklist and migration guide scenarios.

## Evidence

### Governance Normalization

- Previous apply-progress batches were labeled “Strict TDD”.
- Project governance source of truth is `openspec/config.yaml` (`strict_tdd: false`).
- This artifact is normalized to **Standard** mode while preserving prior implementation evidence.

### TDD Cycle Evidence (Historical)

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 1.1 | `config/constructors_test.go` | Unit | ✅ `go test ./config/...` baseline green | ✅ Added failing assertions for new fields/helpers | ✅ `go test ./config -run TestNewJWTConfig\|TestNewRS256ConfigHelpers` | ✅ multiple helper scenarios | ✅ gofmt cleanup |
| 1.2 | `config/constructors_test.go` | Unit | ✅ baseline above | ✅ Helper constructor tests written first | ✅ targeted config tests green | ✅ generated/public/private constructor cases | ✅ gofmt cleanup |
| 1.3 | `config/validate_test.go` | Unit | ✅ baseline above | ✅ Added failing RS256/algorithm matrix tests first | ✅ `go test ./config -run ...Validate...` | ✅ happy + multiple invalid branches | ✅ normalized helper extraction |
| 2.1 | `config/viper/mapping_test.go`, `config/viper/loader_test.go` | Integration | ✅ `go test ./config/...` baseline green | ✅ RS256 mapping/env tests added first | ✅ `go test ./config/viper -run TestMapping\|TestLoadRS256...` | ✅ canonical+legacy+env override paths | ✅ gofmt + focused alias additions |
| 2.2 | `security/keys/keypair_test.go` | Unit | ✅ `go test ./security/...` baseline green | ✅ keypair/resolver tests created before implementation | ✅ `go test ./security/keys -run TestGenerateRSAKeyPair\|TestResolverFromRS256Config` | ✅ generated/public/private + invalid bits/pem/source | ✅ parser helpers extracted |
| 2.3 | `security/jwt/compat_test.go` | Unit | ✅ `go test ./security/...` baseline green | ✅ compat tests wrote new function signature/behaviors first | ✅ `go test ./security/jwt -run TestFromConfigJWT` | ✅ HS256 + RS256 + invalid algorithm/mode | ✅ sentinel cleanup (`errors.New`) |
| 2.4 | `app/application_test.go` | Integration | ✅ `go test ./app/...` baseline green | ✅ helper behavior tests written first | ✅ `go test ./app -run TestUseServerSecurityFromConfig` | ✅ HS256 success, RS256 success, failure no-partial-wiring | ✅ validate-before-wire refactor |
| 3.1 | `config/validate_test.go` | Unit | ✅ baseline preserved | ✅ added matrix before validation code branch changes | ✅ config tests pass | ✅ multi-scenario invalid matrix | ✅ gofmt |
| 3.2 | `config/viper/loader_test.go`, `config/viper/mapping_test.go` | Integration | ✅ baseline preserved | ✅ new RS256 mapping precedence/env tests first | ✅ viper tests pass | ✅ deterministic repeat + precedence checks | ✅ gofmt |
| 3.3 | `security/keys/keypair_test.go` | Unit | ✅ baseline preserved | ✅ tests first for keypair API and errors | ✅ keys tests pass | ✅ success + edge/failure coverage | ✅ helper functions localized |
| 3.4 | `security/jwt/compat_test.go` | Unit | ✅ baseline preserved | ✅ tests first for compat branching | ✅ jwt tests pass | ✅ HS256/RS256 + invalid scenarios | ✅ gofmt |
| 3.5 | `app/application_test.go` | Integration | ✅ baseline preserved | ✅ tests first for wrapper success/failure | ✅ app tests pass | ✅ behavior + no-partial-wiring path | ✅ gofmt |
| 4.1 | `README.md`, `docs/home.md` (docs validation via tests/build) | Docs | ✅ baseline preserved | ✅ N/A (docs) | ✅ full suite/build remained green | ➖ structural docs update | ✅ wording + consistency cleanup |
| 4.2 | `docs/migration/auth-provider-ms-v0.1.0.md` | Docs | ✅ baseline preserved | ✅ N/A (docs) | ✅ full suite/build remained green | ➖ structural docs update | ✅ step ordering cleanup |
| 4.3 | `docs/releases/v0.2.0-checklist.md` | Docs | ✅ baseline preserved | ✅ N/A (docs) | ✅ full suite/build remained green | ➖ structural docs update | ✅ checklist clarity cleanup |
| 4.4 | Full repo (`go test ./...`, `go build ./...`) | Verification | ✅ package-level safety nets completed earlier | ✅ verification commands are acceptance gate | ✅ both commands passed | ➖ single command behavior | ➖ none needed |

## Test Summary

- **Total tests written/updated**: 6 files (new: `security/keys/keypair_test.go`, `security/jwt/compat_test.go`; updated: `config/constructors_test.go`, `config/validate_test.go`, `config/viper/*_test.go`, `app/application_test.go`, `bootstrap_guard_test.go`)
- **Total tests passing**: full suite green via `go test ./...`
- **Layers used**: Unit + Integration
- **Approval tests**: None — this change was additive feature work
- **Pure functions created**: parsing/normalization helpers in config and keypair modules

## Files Changed

- `config/types.go`
- `config/constructors.go`
- `config/constructors_test.go`
- `config/validate.go`
- `config/validate_test.go`
- `config/viper/mapping.go`
- `config/viper/mapping_test.go`
- `config/viper/loader.go`
- `config/viper/loader_test.go`
- `security/keys/keypair.go`
- `security/keys/keypair_test.go`
- `security/jwt/compat.go`
- `security/jwt/compat_test.go`
- `app/application.go`
- `app/application_test.go`
- `README.md`
- `docs/home.md`
- `docs/migration/auth-provider-ms-v0.1.0.md`
- `docs/releases/v0.2.0-checklist.md`
- `bootstrap_guard_test.go`
- `openspec/changes/issue-30-rs256-keypair-security/tasks.md`

## Deviations from Design

None — implementation follows design intent (mode-aware JWT config, provider-agnostic `security/*`, thin app convenience wrapper).

## Issues Found

None.

## Remaining Tasks

None — 16/16 tasks complete.

## Status

Ready for `sdd-verify`.
