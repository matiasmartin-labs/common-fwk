# Tasks: RSA Private Key Accessor

## Phase 1: Foundation — CompatOptions & resolveRS256 helper

- [x] 1.1 `security/jwt/compat.go` — Add `RSAPrivateKey *rsa.PrivateKey` field to `CompatOptions` struct.
- [x] 1.2 `security/jwt/compat.go` — Introduce private `resolveRS256(cfg config.RS256Config) (*rsa.PrivateKey, keys.Resolver, error)` helper; handle `Generated`, `PrivatePEM`, and `PublicPEM` (priv=nil) cases.
- [x] 1.3 `security/jwt/compat.go` — Replace inline RS256 block in `FromConfigJWT` with `resolveRS256` call; populate `CompatOptions.RSAPrivateKey` from the returned key.
  - Acceptance: `CompatOptions.RSAPrivateKey` is non-nil for `Generated`/`PrivatePEM` (spec scenarios: CompatOptions populated for Generated, CompatOptions nil for PublicPEM).

## Phase 2: Core — Application field & accessor

- [x] 2.1 `app/application.go` — Add unexported field `rsaPrivateKey *rsa.PrivateKey` to `Application` struct.
- [x] 2.2 `app/application.go` — In `UseServerSecurityFromConfig`, after compat wiring, assign `a.rsaPrivateKey = compat.RSAPrivateKey`.
- [x] 2.3 `app/application.go` — Add `GetRSAPrivateKey() *rsa.PrivateKey` accessor returning `a.rsaPrivateKey`.
  - Acceptance: Returns non-nil for Generated/PrivatePEM; nil for PublicPEM, direct wiring, or no wiring (spec scenarios §1–§5).

## Phase 3: Testing

- [x] 3.1 `security/jwt/compat_test.go` — Table-driven test for `FromConfigJWT`: RS256 Generated → non-nil, RS256 PrivatePEM → non-nil, RS256 PublicPEM → nil, HS256 → nil.
- [x] 3.2 `app/application_test.go` — Integration test: `UseServerSecurityFromConfig` RS256 Generated → `GetRSAPrivateKey()` non-nil.
- [x] 3.3 `app/application_test.go` — Integration test: `UseServerSecurityFromConfig` RS256 PrivatePEM → `GetRSAPrivateKey()` non-nil.
- [x] 3.4 `app/application_test.go` — Integration test: `UseServerSecurityFromConfig` RS256 PublicPEM → `GetRSAPrivateKey()` nil.
- [x] 3.5 `app/application_test.go` — Unit test: `UseServerSecurity(v)` direct wiring → `GetRSAPrivateKey()` nil.
- [x] 3.6 `app/application_test.go` — Unit test: no security wired → `GetRSAPrivateKey()` nil, no panic.

## Phase 4: Documentation

- [x] 4.1 `app/doc.go` — Add `GetRSAPrivateKey()` to runtime inspection helpers section with nil-safety contract.
- [x] 4.2 `README.md` — Add `GetRSAPrivateKey()` entry in accessors section; note when nil vs non-nil.
- [x] 4.3 `docs/home.md` — Mirror accessor entry from README for doc-sync compliance.
  - Acceptance: All three files describe nil-safety contract consistently (spec scenario: Accessor documented across all surfaces).
