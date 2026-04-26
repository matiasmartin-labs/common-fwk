# Exploration: issue-15-typed-claims

## Current State

`security/claims.Claims` holds the seven standard JWT registered claims (iss, sub, aud, exp, nbf, iat, jti) as typed fields, plus a catch-all `Private map[string]interface{}` for everything else. The `Private` field is populated in `security/jwt/validator.go:claimsFromToken` — any key that is not a registered claim goes into the map.

Key data flow:
1. Raw JWT → `jwt.Validator.Validate()` → `claims.Claims`
2. `security.Validator` interface returns `claims.Claims` (the shared contract)
3. `http/gin/middleware.go` stores the value via `SetClaims(c, key, cl)`
4. Handlers retrieve it via `GetClaims(c, key)` → `claims.Claims`
5. Handlers read domain fields via `cl.Private["email"].(string)` — unsafe, no compile-time guarantee

`GetClaims` / `SetClaims` store and retrieve `claims.Claims` by concrete type assertion (`v.(claims.Claims)`).

## Affected Areas

- `security/claims/claims.go` — `Claims` struct definition; core of the change
- `security/jwt/validator.go` — `claimsFromToken()` populates `Private`; `Validator` interface returns `claims.Claims`
- `security/validator.go` — top-level `Validator` interface returns `claims.Claims`
- `http/gin/context.go` — `SetClaims` / `GetClaims` typed on `claims.Claims`
- `http/gin/middleware.go` — stores claims; signature may need updating
- `security/claims/claims_test.go` — must cover new access patterns

## Approaches

### Option A — Add common OIDC fields directly to Claims struct

Add `Email`, `Name`, `Roles`, `Picture` (and similar standard OIDC fields) as first-class fields on `Claims`, populated from `Private` inside `claimsFromToken`.

```go
type Claims struct {
    // standard JWT
    Issuer    string    `json:"iss,omitempty"`
    Subject   string    `json:"sub,omitempty"`
    // ...
    // OIDC profile fields
    Email     string    `json:"email,omitempty"`
    Name      string    `json:"name,omitempty"`
    Picture   string    `json:"picture,omitempty"`
    Roles     []string  `json:"roles,omitempty"` // or []string based on OIDC spec
    // remaining unknown fields
    Private   map[string]interface{} `json:"-"`
}
```

- **Pros**: Zero boilerplate for consumers; fully compile-time safe; backward compatible (existing code reading `Private` still works for non-OIDC keys); no generics; simple to understand
- **Cons**: `Claims` accumulates fields over time; not every issuer emits all fields (empty strings vs absent); `claims` package becomes opinionated about OIDC — mixing concern levels (JWT spec + OIDC profile)
- **Effort**: Low

---

### Option B — Generic `TypedClaims[T any]` wrapper

Introduce a wrapper that carries both standard claims and a typed payload.

```go
type TypedClaims[T any] struct {
    Claims              // embedded standard claims
    Payload T
}
```

Consumers define their own struct, then use a helper to bind:
```go
type AppClaims struct {
    Email   string   `json:"email"`
    Roles   []string `json:"roles"`
}

tc, err := claims.Bind[AppClaims](cl) // extracts from Private map into T
```

The `Validator` interface stays returning `claims.Claims`; `Bind` is called post-validation in the handler or middleware layer.

- **Pros**: Maximum type safety; consumers own their schema; no accumulation of fields in the core struct; works for ANY domain
- **Cons**: Requires Go 1.18+ generics (satisfied — module is go 1.25); `Bind` adds a step and an extra error path; `GetClaims` in gin/context.go returns `claims.Claims`, so `Bind` must be called after retrieval; `TypedClaims[T]` cannot be stored via existing `SetClaims` without also adding a generic `SetTypedClaims[T]` — creates an API surface split
- **Effort**: Medium

---

### Option C — Typed accessor helpers for the Private map

Keep the struct unchanged, add type-safe getter functions:

```go
func GetString(c Claims, key string) (string, bool)
func GetStrings(c Claims, key string) ([]string, bool)
// convenience aliases
func Email(c Claims) (string, bool)   { return GetString(c, "email") }
func Name(c Claims) (string, bool)    { return GetString(c, "name") }
```

- **Pros**: Zero struct changes; completely backward compatible; no generics needed; easy to add new accessors
- **Cons**: Still not compile-time safe — mistyped key names are silent failures; the string-keyed map remains the actual data store; just hides the map, doesn't eliminate it; tests must cover all accessors; adds API surface without real safety gain
- **Effort**: Low

---

## Recommendation

**Option A** — add common OIDC fields directly to `Claims`, keeping `Private` for unknown extras.

**Justification**:
1. The fields in question (`email`, `name`, `picture`, `roles`) are standardized OIDC profile claims (RFC 7519 + OIDC Core 1.0). They are not arbitrary domain fields — they belong in a foundational claims type.
2. Zero consumer friction: existing code continues to work. New code accesses `cl.Email` with no migration effort.
3. The `jwt.claimsFromToken` already maps standard JWT fields explicitly — extending this for well-known OIDC fields is idiomatic and consistent with the existing pattern.
4. Option B (generics) adds complexity that pays off only when consumers need truly custom schemas. For the stated use case (auth-provider-ms needing email/name/roles/picture), Option A is sufficient and simpler.
5. Option C is strictly inferior: it hides but doesn't eliminate the map; type errors are deferred to runtime.

**Scope boundary**: Only well-known OIDC profile fields go into the struct. Truly proprietary claims (e.g., tenant-specific metadata) stay in `Private`. This keeps `Claims` principled rather than becoming a dumping ground.

## Risks

- `Roles` claim name may differ across issuers (e.g., Auth0 uses `https://example.com/roles` as a namespaced claim). Should document that `Roles` maps to the plain `roles` key; namespaced variants remain in `Private`.
- Adding fields to `Claims` is a minor API expansion — not breaking, but consumers comparing structs with `==` will now see more fields (unlikely in practice).
- `email_verified` and other common OIDC fields may be requested as follow-ups; the struct can grow incrementally but should be documented.

## Ready for Proposal

Yes — Option A is well-scoped, low-risk, and directly addresses the issue. The orchestrator should proceed to `sdd-propose`, then `sdd-spec`, then `sdd-tasks`.
