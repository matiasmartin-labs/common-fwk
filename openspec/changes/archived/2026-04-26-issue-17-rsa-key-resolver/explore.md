# Exploration: RSA Key Resolver for RS256 Support

## Current State

### `security/keys` package

- `Key` struct has three fields: `ID string`, `Method string`, `Verify any`
  - `Verify` is already typed as `any` — not `[]byte` as the issue suggested. This is a **key finding**: the struct already accommodates arbitrary verify material.
- `Resolver` interface: `Resolve(ctx context.Context, kid string) (Key, error)`
- `NewStaticResolver(defaultKey *Key, byID map[string]Key) Resolver` — in-memory map lookup with optional default fallback.
- No existing RSA-specific resolver or constructor.

### `security/jwt` package

- `validator.Validate` calls `v.options.Resolver.Resolve(ctx, kid)` to get a `Key`.
- Passes `key.Verify` directly to `golang-jwt/v5`'s keyfunc:
  ```go
  return key.Verify, nil
  ```
- `golang-jwt/v5` dispatches to the signing method's `Verify(signingString, signature string, key interface{}) error`. For RS256 it expects `*rsa.PublicKey`; for HS256 it expects `[]byte`.
- **No changes are needed in `validator.go`** — it already passes `key.Verify` as `any` to the library, which handles type-dispatch internally.
- `Options.Methods` defaults to `["HS256"]` — callers must explicitly add `"RS256"`.

### Test helpers

- `mustSignToken` uses `[]byte` secret for HS256 signing.
- Tests are `package jwt` (white-box). New RS256 tests can live in the same file or a separate `_test.go`.

## Affected Areas

- `security/keys/types.go` — add `NewRSAKey` constructor helper and/or document RSA usage; `Verify any` field already works.
- `security/keys/resolver.go` — add `NewRSAResolver(privateKey *rsa.PrivateKey, keyID string) Resolver`.
- `security/jwt/validator_test.go` — add RS256 test cases (valid, invalid sig, expired).
- `security/jwt/options.go` — no changes needed; callers supply `Methods: []string{"RS256"}`.
- `security/jwt/validator.go` — **no changes needed**.

## Approaches

### Approach 1 — Add `NewRSAResolver` only (minimal)

Add a single constructor in `security/keys/resolver.go` (or a new `rsa_resolver.go`) that wraps `NewStaticResolver` with the public key in `Key.Verify`:

```go
func NewRSAResolver(priv *rsa.PrivateKey, keyID string) Resolver {
    k := Key{ID: keyID, Method: "RS256", Verify: &priv.PublicKey}
    return NewStaticResolver(&k, nil)
}
```

- **Pros**: Zero changes to `Key`, `Resolver` interface, or `validator.go`; minimal surface area; backward compatible; matches `golang-jwt` expectations directly.
- **Cons**: Caller must hold the private key only to extract the public key — slightly odd API if they only have the public key (e.g., a verify-only service).
- **Effort**: Low

### Approach 2 — Add both `NewRSAResolver` and `NewRSAPublicKeyResolver`

Two constructors: one taking `*rsa.PrivateKey` (derives public key), one taking `*rsa.PublicKey` directly.

- **Pros**: More flexible; covers verify-only consumers that never have the private key.
- **Cons**: Slightly more surface area; still trivially small.
- **Effort**: Low

### Approach 3 — Add `PublicKey crypto.PublicKey` field to `Key` struct

As the issue proposes: add a new field alongside `Verify`, and modify `validator.go` to choose the right field.

- **Pros**: Explicit typed field; self-documenting struct.
- **Cons**: `Verify` is already `any` and already works; adding a parallel field creates ambiguity (which one wins?), requires modifying `validator.go`, and adds complexity for no functional gain.
- **Effort**: Medium, **higher risk** of breakage.

## Recommendation

**Approach 2** — add both `NewRSAResolver(priv *rsa.PrivateKey, keyID string) Resolver` and `NewRSAPublicKeyResolver(pub *rsa.PublicKey, keyID string) Resolver`. Both are trivial wrappers around `NewStaticResolver`. Zero changes to `Key`, `Resolver`, `validator.go`, or `Options`. The `Verify any` field already works because `golang-jwt/v5` dispatches by type internally.

The issue's proposal to add `PublicKey crypto.PublicKey` to the `Key` struct is **not necessary** — `Verify any` already holds it.

## Risks

- `golang-jwt/v5` RS256 keyfunc receives `*rsa.PublicKey` — if someone accidentally passes `*rsa.PrivateKey` instead, the library will return an error at verify time. Constructor should use `&priv.PublicKey` explicitly.
- `Options.Methods` defaults to `["HS256"]`; RS256 callers **must** set `Methods: []string{"RS256"}` or `["HS256", "RS256"]`. This is a documentation/example risk rather than a code risk.
- Existing tests use `package jwt` (internal); new RS256 tests can reuse `mustSignToken` with `jwtlib.SigningMethodRS256` and an `*rsa.PrivateKey`.

## Ready for Proposal

**Yes.** The scope is well-defined, the implementation path is clear, and there is zero risk to existing HS256 behavior.
