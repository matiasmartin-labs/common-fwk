# Tasks: RS256 Keypair Security Bootstrap

## Phase 1: Config Foundation

- [x] 1.1 Extend `config/types.go` with `JWTConfig.Algorithm` (default `HS256`) and `RS256Config` (`KeySource`, `KeyID`, `PublicKeyPEM`, `PrivateKeyPEM`) while preserving existing fields.
- [x] 1.2 Update `config/constructors.go` so `NewJWTConfig(secret, issuer, ttlMinutes)` remains backward-compatible and add focused RS256 helper/default constructors used by tests.
- [x] 1.3 Implement conditional JWT validation in `config/validate.go` for HS256 vs RS256 requirements, including deterministic errors for unsupported algorithm/mode combinations.

## Phase 2: Adapter + Security Core Wiring

- [x] 2.1 Add RS256 canonical key mapping in `config/viper/mapping.go` and deterministic compatibility alias behavior in `config/viper/loader.go` for file+env inputs.
- [x] 2.2 Create `security/keys/keypair.go` with in-memory RSA keypair generation and resolver bootstrap (`generated`, `public-pem`, `private-pem`) without provider coupling.
- [x] 2.3 Update `security/jwt/compat.go` so `FromConfigJWT` branches by algorithm and returns algorithm-constrained `Options` (`Methods`, resolver, issuer, TTL) for HS256 and RS256.
- [x] 2.4 Add `UseServerSecurityFromConfig()` to `app/application.go` as a thin wrapper that builds validator from config and delegates to existing `UseServerSecurity` with wrapped errors.

## Phase 3: Deterministic Tests and Invalid-Mode Coverage

- [x] 3.1 Expand `config/validate_test.go` with table-driven HS256/RS256 matrices, including invalid algorithm, missing `secret`, missing `key_id`, and missing PEM scenarios.
- [x] 3.2 Update `config/viper/*_test.go` to verify deterministic RS256 mapping precedence, legacy alias compatibility, env overrides, and typed failure paths.
- [x] 3.3 Create `security/keys/keypair_test.go` for valid generation/retrieval plus invalid input cases (bad bits, malformed PEM, missing key) with assertable error classification.
- [x] 3.4 Add `security/jwt/compat_test.go` validating `FromConfigJWT` HS256/RS256 resolver wiring, method allowlist correctness, and deterministic invalid-mode failures.
- [x] 3.5 Update `app/application_test.go` to cover `UseServerSecurityFromConfig()` success ordering and no-partial-wiring behavior on invalid configuration.

## Phase 4: Documentation and Verification

- [x] 4.1 Update `README.md` and `docs/home.md` with configuration examples for HS256 default and RS256 bootstrap semantics.
- [x] 4.2 Update `docs/migration/auth-provider-ms-v0.1.0.md` with an executable HS256→RS256 migration sequence and verification checklist.
- [x] 4.3 Update `docs/releases/v0.2.0-checklist.md` with release checks for HS256 compatibility and RS256 bootstrap readiness.
- [x] 4.4 Run `go test ./...` and `go build ./...` and record any follow-up fixes required to keep the change releasable, including doc-contract assertions for release checklist and migration guide scenarios.
