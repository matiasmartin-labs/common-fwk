## Exploration: issue-52-rsa-private-key-accessor

### Current State

#### `app.Application` struct

Defined in `app/application.go`, the struct holds:

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
}
```

**There is no field for storing an RSA private key.** The private key is consumed during security wiring and not retained.

#### How the private key is set (flow)

1. `UseServerSecurityFromConfig()` in `app/application.go`:
   - Calls `securityjwt.FromConfigJWT(cfg.Security.Auth.JWT)` → returns `CompatOptions`
   - Calls `securityjwt.NewValidator(compat.Options)` → creates validator
   - Calls `a.UseServerSecurity(validator)` → stores only the `security.Validator` interface

2. `securityjwt.FromConfigJWT` in `security/jwt/compat.go`:
   - For RS256: calls `keys.ResolverFromRS256Config(cfg.RS256)` → builds a `Resolver`
   - The `Resolver` is embedded in `Options.Resolver` — **the private key is never surfaced**

3. `keys.ResolverFromRS256Config` in `security/keys/keypair.go`:
   - For `KeySourceGenerated`: generates a new `*rsa.PrivateKey` → passes to `NewRSAResolver(priv, ...)` → **key is buried inside a `staticResolver` struct**
   - For `KeySourcePrivatePEM`: parses PEM → passes to `NewRSAResolver(priv, ...)` → **key is buried inside a `staticResolver` struct**
   - For `KeySourcePublicPEM`: only a `*rsa.PublicKey` is used; **no private key available at all**

4. `NewRSAResolver` in `security/keys/rsa.go`:
   - Extracts `&privateKey.PublicKey` for the `Key.Verify` field
   - The original `*rsa.PrivateKey` pointer is **dropped** — only the public key survives in the resolver

**Root cause**: The private key is discarded immediately after the `Resolver` is constructed. Nothing in the call chain between keypair loading and `Application` retains it.

#### Existing accessor patterns

The following read-only accessors exist on `Application`:
- `GetConfig() config.Config` — defensive snapshot (deep copy)
- `GetSecurityValidator() security.Validator` — returns interface as-is, nil when not wired
- `IsSecurityReady() bool` — boolean readiness flag

Pattern:
- No guard/error return — nil is the zero value for pointer and interface returns
- Nil-safe: callers check for nil themselves
- No locking: `Application` is not designed for concurrent mutation after bootstrap

#### Issue #50 patterns (JWKS / public key exposure)

`NewRSAPublicKeyResolver` in `security/keys/rsa.go` follows the same "only public key" model. No existing JWKS endpoint or public-key accessor exists on `Application`. This change is independent of #50.

---

### Affected Areas

- `app/application.go` — add `rsaPrivateKey *rsa.PrivateKey` field, capture key during wiring, add `GetRSAPrivateKey()` accessor
- `app/application_test.go` — add tests for `GetRSAPrivateKey()` covering: nil when not wired (HS256), nil when public-key-only RS256, non-nil when private-PEM or generated RS256
- `security/jwt/compat.go` — `FromConfigJWT` must surface the `*rsa.PrivateKey` so the `app` layer can capture it; `CompatOptions` should include the key
- `security/keys/keypair.go` — `ResolverFromRS256Config` must return the key alongside the resolver (or a different struct/return signature)
- Possibly: `app/doc.go`, `README.md`, `docs/home.md` — documentation sync tests enforce that accessor signatures appear in docs

---

### Approaches

#### Approach 1: Extend `CompatOptions` to carry the private key

Extend `CompatOptions` in `security/jwt/compat.go`:
```go
type CompatOptions struct {
    Options        Options
    TokenTTL       time.Duration
    RSAPrivateKey  *rsa.PrivateKey // non-nil only when RS256 with private key source
}
```
`FromConfigJWT` sets `RSAPrivateKey` for PrivatePEM and Generated sources.  
`UseServerSecurityFromConfig` stores `compat.RSAPrivateKey` to a new `Application.rsaPrivateKey` field.

- **Pros**: Minimal blast radius; all changes stay in the `app` + `jwt` packages. `CompatOptions` is the right abstraction to carry issuer-side data (it already carries `TokenTTL`). Clean layering.
- **Cons**: `CompatOptions` grows one more field; `security/keys` package is not changed.
- **Effort**: Low

#### Approach 2: `ResolverFromRS256Config` returns a richer result struct

Define a `RS256ResolverResult{Resolver, PrivateKey}` in `security/keys` and return from `ResolverFromRS256Config`.

- **Pros**: Private key is surfaced at the lowest layer, allows future callers beyond `app`.
- **Cons**: Wider API surface change in `security/keys`; all callers of `ResolverFromRS256Config` need updating. Overkill for this issue.
- **Effort**: Medium

#### Approach 3: New `WithRSAPrivateKey` bootstrap method

Add a separate `UseServerSecurityWithRSA(v security.Validator, priv *rsa.PrivateKey)` method to `Application`, leaving `UseServerSecurity` unchanged.

- **Pros**: Backward-compatible; callers who manage keys separately can supply both.
- **Cons**: Does not solve the case where `UseServerSecurityFromConfig()` generates the key — caller cannot obtain it externally. Two paths to set security diverge.
- **Effort**: Medium; only shifts the responsibility to caller.

---

### Recommendation

**Approach 1 — Extend `CompatOptions`** is recommended.

- `CompatOptions` already exists as the bridge between config-land and app-land for issuer concerns (`TokenTTL`). Adding `RSAPrivateKey *rsa.PrivateKey` is consistent with that purpose.
- Changes are confined to three files: `security/jwt/compat.go`, `app/application.go`, `app/application_test.go`.
- The accessor follows the existing nil-return pattern (`GetSecurityValidator` returns nil when not set).
- The `KeySourcePublicPEM` case correctly returns `nil` because there is no private key — nil-return is the right contract.
- No changes to `security/keys` are needed.

**Implementation sketch**:
```go
// security/jwt/compat.go — CompatOptions
type CompatOptions struct {
    Options       Options
    TokenTTL      time.Duration
    RSAPrivateKey *rsa.PrivateKey // non-nil only for private-key RS256 sources
}

// app/application.go — Application struct
type Application struct {
    // ... existing fields ...
    rsaPrivateKey  *rsa.PrivateKey
}

// app/application.go — UseServerSecurityFromConfig
func (a *Application) UseServerSecurityFromConfig() (*Application, error) {
    // ... existing logic ...
    a.rsaPrivateKey = compat.RSAPrivateKey
    // ...
}

// app/application.go — GetRSAPrivateKey
// GetRSAPrivateKey returns the RSA private key when security was wired with an
// RS256 private-key source. Returns nil if security was not wired, or if only a
// public key was configured (e.g. public-PEM or HS256).
func (a *Application) GetRSAPrivateKey() *rsa.PrivateKey {
    return a.rsaPrivateKey
}
```

---

### Risks

1. **Public-key-only RS256**: When `KeySourcePublicPEM` is used, there is no private key and `GetRSAPrivateKey()` must return nil. This is by design and must be documented clearly to avoid caller confusion.
2. **Generated key not reproducible**: When `KeySourceGenerated` is used, each call to `UseServerSecurityFromConfig` regenerates the key. Applications relying on the generated key for token issuance across restarts will have token invalidation. This is a pre-existing concern, not introduced by this accessor.
3. **Thread safety**: `Application` is not designed for concurrent mutation; bootstrap is sequential by convention. The accessor is read-only post-bootstrap, so no locking is needed. This must be documented (matches existing pattern for `GetSecurityValidator`).
4. **`UseServerSecurity(validator)` direct path**: When the caller uses `UseServerSecurity` directly (bypassing `UseServerSecurityFromConfig`), `rsaPrivateKey` will always be nil. The accessor correctly reflects this — callers wiring security manually are responsible for their own key management.
5. **Documentation sync tests**: `TestDocumentation_AccessorContractSynchronization` checks for accessor signatures in `doc.go`, `README.md`, and `docs/home.md`. A new accessor will require documentation updates to keep that test passing.

---

### Ready for Proposal

Yes. The exploration is complete. The recommended approach (Approach 1) is unambiguous, low-risk, and consistent with existing patterns. The proposal phase can proceed with these constraints:

- Extend `CompatOptions.RSAPrivateKey`
- Add `Application.rsaPrivateKey` field
- Capture key in `UseServerSecurityFromConfig` after `FromConfigJWT`
- Add `GetRSAPrivateKey() *rsa.PrivateKey` accessor
- Add tests in `app/application_test.go`
- Update doc files to satisfy existing documentation sync tests
