# Tasks: Issue #4 Security Core JWT Validation

## Phase 1: Foundation Contracts (claims + keys)

- [x] 1.1 Create `security/claims/doc.go` and `security/claims/claims.go` with standard claims model and `aud` normalization helper.
  Completion: package compiles and supports `iss/sub/aud/exp/nbf/iat/jti` plus private claims without validation side effects.
- [x] 1.2 Create `security/keys/doc.go`, `security/keys/types.go`, and `security/keys/resolver.go` with `Key` and `Resolver` contracts plus deterministic in-memory resolver.
  Completion: resolver returns `kid` match, default key fallback, and categorized miss error without network access.

## Phase 2: Validator Core + Error Contracts

- [x] 2.1 Create `security/jwt/doc.go` and `security/jwt/errors.go` with sentinel taxonomy and typed `ValidationError` wrapper.
  Completion: malformed/signature/issuer/audience/method/expired/not-yet-valid/key-resolution categories are exported and `errors.Is/As` compatible.
- [x] 2.2 Create `security/jwt/options.go` with `Methods`, `Issuer`, `Audience`, `Now`, and `Resolver` options, including safe defaults.
  Completion: options can be constructed without globals and default `Now` is deterministic when overridden in tests.
- [x] 2.3 Create `security/jwt/validator.go` implementing ordered flow: parse → method gate → key resolve → signature verify → claims checks.
  Completion: `Validate(ctx, raw)` returns normalized claims on success and stage-appropriate wrapped errors on failure.

## Phase 3: Compatibility + Documentation

- [x] 3.1 Create `security/jwt/compat.go` mapping `config.JWTConfig` (`Secret`, `Issuer`, `TTLMinutes`) into validator options without adding runtime coupling to `config` validation.
  Completion: mapping preserves backward-compatible JWT field semantics and keeps `config` package configuration-only.
- [x] 3.2 Update `README.md` with security-core usage examples and explicit non-goals (no Gin middleware/JWKS adapter scope).
  Completion: README documents package boundaries, validator setup, and compatibility mapping behavior.

## Phase 4: Tests and Verification

- [x] 4.1 Add `security/claims/claims_test.go` table tests for `aud` string/array normalization and optional-claim parsing.
  Completion: tests cover equivalent audience encodings and missing optional claims scenarios.
- [x] 4.2 Add `security/keys/resolver_test.go` for `kid` hit/miss and default-key behavior.
  Completion: resolver tests assert deterministic outcomes and key-resolution error category.
- [x] 4.3 Add `security/jwt/validator_test.go` deterministic tests (fixed clock + fake resolver) for happy path and each failure category.
  Completion: tests assert method rejection, invalid issuer/audience, expired, not-before, malformed token, invalid signature, and key-resolution failures.
- [x] 4.4 Add contract tests (can be in `security/jwt/validator_test.go`) that wrap returned errors and assert `errors.Is`/`errors.As` stability.
  Completion: wrapped validation errors remain assertable by sentinel and typed wrapper.
- [x] 4.5 Run validation commands: `go test ./...` then `go build ./...`.
  Completion: both commands succeed with no failing packages.
