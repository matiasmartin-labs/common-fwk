## Exploration: issue-53-rsa-public-key-accessor

### Current State

#### `app.Application` struct (post-issue-52)

Defined in `app/application.go`:

```go
type Application struct {
    cfg            config.Config
    server         http.Server
    handler        *gin.Engine
    validator      security.Validator
    loggerRegistry logging.Registry
    logOutput      io.Writer
    serverReady    bool
    securityReady  bool
    rsaPrivateKey  *rsa.PrivateKey   // added by issue-52
}
```

Issue #52 already added `rsaPrivateKey` and `GetRSAPrivateKey()`. The public key accessor is its direct companion.

#### Where the RSA public key lives

The public key is reachable from the **existing** `rsaPrivateKey` field as `&a.rsaPrivateKey.PublicKey` (for Generated and PrivatePEM sources). However, when `KeySourcePublicPEM` is used, `rsaPrivateKey` is nil — only the public key is available, and it is **not** currently retained anywhere on `Application`.

#### How security wiring works (full chain)

1. `UseServerSecurityFromConfig()` calls `securityjwt.FromConfigJWT(cfg.Security.Auth.JWT)` → returns `CompatOptions`
2. `CompatOptions` currently carries:
   - `Options` — validator options including the `keys.Resolver`
   - `TokenTTL` — duration
   - `RSAPrivateKey *rsa.PrivateKey` — non-nil for Generated/PrivatePEM, nil for PublicPEM
3. For RS256, `resolveRS256` in `security/jwt/compat.go` calls:
   - `keys.NewRSAResolver(priv, cfg.KeyID)` for Generated/PrivatePEM — builds a `staticResolver` with `Key{ID: keyID, Verify: &priv.PublicKey}`
   - `keys.ResolverFromRS256Config(cfg)` → `keys.NewRSAPublicKeyResolver(pub, cfg.KeyID)` for PublicPEM
4. **The KeyID** (from `cfg.RS256.KeyID`) is stored in the `Key` inside the `staticResolver` but not surfaced to `Application`
5. **The public key** is stored as `Key.Verify` (of type `*rsa.PublicKey`) inside the resolver — not surfaced to `Application`

#### Key ID storage

`config.RS256Config.KeyID` is used to construct the `keys.Key.ID` inside the resolver. It is **not** retained by `Application` directly. The key ID is required for all three key sources (Generated, PrivatePEM, PublicPEM) — `ErrRS256KeyIDRequired` is enforced.

#### Issue #52 `GetRSAPrivateKey()` implementation

```go
// Application struct field
rsaPrivateKey *rsa.PrivateKey

// Accessor
func (a *Application) GetRSAPrivateKey() *rsa.PrivateKey {
    return a.rsaPrivateKey
}

// Captured in UseServerSecurityFromConfig
a.rsaPrivateKey = compat.RSAPrivateKey
```

Pattern: zero-overhead, nil-safe, no locking, no error return.

#### `CompatOptions` structure

```go
type CompatOptions struct {
    Options       Options
    TokenTTL      time.Duration
    RSAPrivateKey *rsa.PrivateKey // non-nil only for RS256 Generated/PrivatePEM sources
}
```

#### Test patterns in `app/application_test.go`

- Tests use `t.Parallel()` at the function level and subtests use `t.Parallel()` too
- Table-driven structure with `setup func(t *testing.T) *Application` + `wantNonNil bool`
- Helper `mustNotPanic(t, name, fn)` ensures accessors never panic regardless of state
- Helper `mustRSAPrivatePEM(t)` / `mustRSAPublicPEM(t)` generate keys inline
- Helper `rs256AppConfig(rs256Cfg)` builds config for a given RS256 config
- All three key sources are exercised: Generated, PrivatePEM, PublicPEM
- `UseServerSecurity` direct path (no RSA) returns nil
- "no security wired" returns nil without panic
- `TestDocumentation_AccessorContractSynchronization` enforces accessor signatures exist in `doc.go`, `README.md`, `docs/home.md`

---

### Affected Areas

- `app/application.go` — add `rsaPublicKey *rsa.PublicKey` and `rsaKeyID string` fields; capture in `UseServerSecurityFromConfig`; add `GetRSAPublicKey()` and `GetRSAKeyID()` accessors
- `app/application_test.go` — add `TestGetRSAPublicKey` and `TestGetRSAKeyID` table-driven tests covering all three RS256 sources plus HS256 and unwired paths
- `security/jwt/compat.go` — extend `CompatOptions` with `RSAPublicKey *rsa.PublicKey` and `RSAKeyID string`; populate them in `FromConfigJWT` / `resolveRS256`
- `app/doc.go` — add accessor signatures and lifecycle contract for `GetRSAPublicKey()` and `GetRSAKeyID()` (required by documentation sync test)
- `README.md` and `docs/home.md` — update accessor tables to include new signatures (required by documentation sync test)

---

### Approaches

#### Approach 1: Extend `CompatOptions` (mirrors issue-52 pattern)

Add `RSAPublicKey *rsa.PublicKey` and `RSAKeyID string` to `CompatOptions`. Populate them in `resolveRS256`. Capture in `UseServerSecurityFromConfig`.

- **Pros**: Exact same pattern as issue-52 (`RSAPrivateKey`). Minimal blast radius. Changes confined to `security/jwt/compat.go` and `app/`. All three key sources can supply the public key.
- **Cons**: `CompatOptions` grows two more fields.
- **Effort**: Low

#### Approach 2: Derive public key from `rsaPrivateKey` at accessor call time

For Generated/PrivatePEM: return `&a.rsaPrivateKey.PublicKey`. For PublicPEM or when no private key: need a separate field anyway. This is a partial shortcut.

- **Pros**: No change to `CompatOptions` for the Generated/PrivatePEM case.
- **Cons**: Requires a separate `rsaPublicKey` field anyway for the PublicPEM case; inconsistent with how `rsaPrivateKey` is handled. Two code paths for one accessor is harder to reason about.
- **Effort**: Low-to-Medium (still requires `CompatOptions` change for PublicPEM)

#### Approach 3: Parse public key from resolver at accessor call time

Cast `resolver.Resolve(ctx, keyID).Verify` to `*rsa.PublicKey`. No new fields on `Application`.

- **Pros**: No struct changes.
- **Cons**: Requires a `context.Context` or stored `kid`; the resolver is not stored on `Application` (only the `validator` is); would require surfacing the resolver itself. Violates the accessor pattern (should be O(1) field read, not a resolve call).
- **Effort**: High (invasive)

---

### Recommendation

**Approach 1 — Extend `CompatOptions`** following the exact issue-52 pattern.

Add to `CompatOptions`:
```go
type CompatOptions struct {
    Options       Options
    TokenTTL      time.Duration
    RSAPrivateKey *rsa.PrivateKey // non-nil for Generated/PrivatePEM
    RSAPublicKey  *rsa.PublicKey  // non-nil for all RS256 sources (Generated/PrivatePEM/PublicPEM)
    RSAKeyID      string          // non-empty for all RS256 sources
}
```

In `resolveRS256`:
- Generated: `priv` → `RSAPublicKey = &priv.PublicKey`, `RSAKeyID = cfg.KeyID`
- PrivatePEM: `priv` → `RSAPublicKey = &priv.PublicKey`, `RSAKeyID = cfg.KeyID`
- PublicPEM: `pub` from parse → `RSAPublicKey = pub`, `RSAKeyID = cfg.KeyID`

In `Application`:
```go
rsaPublicKey *rsa.PublicKey
rsaKeyID     string
```

Accessors:
```go
// GetRSAPublicKey returns the RSA public key used for RS256 token verification.
// Returns nil when security was not wired, or when the algorithm is not RS256.
func (a *Application) GetRSAPublicKey() *rsa.PublicKey {
    return a.rsaPublicKey
}

// GetRSAKeyID returns the key ID associated with the RS256 key.
// Returns empty string when security was not wired or algorithm is not RS256.
func (a *Application) GetRSAKeyID() string {
    return a.rsaKeyID
}
```

---

### Risks

1. **PublicPEM returns nil private key but non-nil public key**: This is by design and the correct semantic — document clearly. `GetRSAPublicKey()` returns non-nil for ALL three RS256 key sources; `GetRSAPrivateKey()` returns nil for PublicPEM only.
2. **KeyID is always required for RS256**: `ErrRS256KeyIDRequired` is enforced upstream; `GetRSAKeyID()` will never return empty when properly wired with RS256.
3. **Documentation sync test**: `TestDocumentation_AccessorContractSynchronization` requires accessor signatures in `doc.go`, `README.md`, `docs/home.md`. Both new accessor signatures must be added.
4. **HS256 path**: `GetRSAPublicKey()` returns nil and `GetRSAKeyID()` returns `""` when algorithm is HS256 — correct by design.
5. **Direct wiring via `UseServerSecurity(v)`**: Both accessors return nil/empty — consistent with `GetRSAPrivateKey()` behavior.

---

### Ready for Proposal

Yes. The approach is clear, low-risk, and directly mirrors issue-52. The proposal phase can proceed.
