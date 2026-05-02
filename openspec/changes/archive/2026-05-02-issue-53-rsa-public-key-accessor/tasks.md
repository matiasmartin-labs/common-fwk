# Tasks: RSA Public Key and Key ID Accessors

> TDD mode active. RED → GREEN → REFACTOR order is mandatory.

## Phase 1: RED — Write Failing Tests

- [x] 1.1 `app/application_test.go` — Add `TestGetRSAPublicKey` table-driven test covering 5 scenarios: Generated, PrivatePEM, PublicPEM, HS256, unwired. Must compile but fail (methods don't exist yet).
- [x] 1.2 `app/application_test.go` — Add `TestGetRSAKeyID` table-driven test covering 5 scenarios: Generated, PrivatePEM, PublicPEM, HS256, unwired. Must compile but fail.
- [x] 1.3 Run `go test ./app/...` — confirmed red (compile error).

## Phase 2: GREEN — Implement Production Code

- [x] 2.1 `security/jwt/compat.go` — Add `RSAPublicKey *rsa.PublicKey` and `RSAKeyID string` fields to `CompatOptions`.
- [x] 2.2 `security/jwt/compat.go` — Update `resolveRS256` signature to `(*rsa.PrivateKey, *rsa.PublicKey, keys.Resolver, error)`; populate public key in all branches; add `parseRS256PublicPEM` for PublicPEM source; set `RSAKeyID = cfg.RS256.KeyID` in `CompatOptions` assembly.
- [x] 2.3 `app/application.go` — Add private fields `rsaPublicKey *rsa.PublicKey` and `rsaKeyID string`.
- [x] 2.4 `app/application.go` — In `UseServerSecurityFromConfig`, capture `compat.RSAPublicKey` → `a.rsaPublicKey` and `compat.RSAKeyID` → `a.rsaKeyID`.
- [x] 2.5 `app/application.go` — Add `GetRSAPublicKey() *rsa.PublicKey { return a.rsaPublicKey }` and `GetRSAKeyID() string { return a.rsaKeyID }`.
- [x] 2.6 Run `go test ./...` — confirmed all tests pass (green).

## Phase 3: Documentation

- [x] 3.1 `app/doc.go` — Add `GetRSAPublicKey() *rsa.PublicKey` and `GetRSAKeyID() string` to the accessor contract block.
- [x] 3.2 Run `go test ./app/...` — confirmed passes.
- [x] 3.3 `README.md` — Add accessor rows.
- [x] 3.4 `docs/home.md` — Add same accessor rows.

## Phase 4: Final Verification

- [x] 4.1 Run `go test ./...` — full suite green, no regressions.
