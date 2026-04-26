# Design: Typed OIDC Claims on claims.Claims

## Technical Approach

Add `Email`, `Name`, `Picture string` and `Roles []string` as first-class fields on `claims.Claims`. Populate them in `claimsFromToken()` alongside existing standard claims, using the same `mapClaims["key"].(type)` pattern already used for `iss`, `sub`, `jti`. Exclude extracted keys from `Private` by adding them to the switch-case skip list. No interface or package-level API changes.

## Architecture Decisions

| Decision | Choice | Alternatives | Rationale |
|---|---|---|---|
| Field location | `claims.Claims` struct fields | Wrapper type / separate struct | Follows existing flat struct pattern; no breaking change |
| `Roles` key | Only bare `roles` | Also `realm_access.roles`, namespaced keys | Namespaced variants are IdP-specific; they stay in `Private` to avoid coupling |
| `Roles` nil vs empty | `nil` when key absent, slice when present | Always allocate | Consistent with `Private` map — only populated when data exists |
| Extraction site | `claimsFromToken()` in `jwt/validator.go` | `Claims` methods, middleware | Single extraction point; mirrors how all other fields are populated |

## Data Flow

```
JWT token (raw string)
  └─► claimsFromToken(token *jwtlib.Token)
        │
        ├─► mapClaims["email"]  ──► Claims.Email   (string)
        ├─► mapClaims["name"]   ──► Claims.Name    (string)
        ├─► mapClaims["picture"]──► Claims.Picture (string)
        ├─► mapClaims["roles"]  ──► Claims.Roles   ([]string, dual type assertion)
        │
        └─► remaining keys ──► Claims.Private (map[string]interface{})
              (skips: iss, sub, aud, exp, nbf, iat, jti, email, name, picture, roles)
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `security/claims/claims.go` | Modify | Add `Email`, `Name`, `Picture string` and `Roles []string` fields to `Claims` struct |
| `security/jwt/validator.go` | Modify | Extract OIDC fields in `claimsFromToken()`; extend skip-list in Private loop |
| `security/jwt/validator_test.go` | Modify | Add table-driven cases for tokens with OIDC fields; test `[]interface{}` and `[]string` roles variants; test absent fields |

## Interfaces / Contracts

```go
// security/claims/claims.go
type Claims struct {
    Issuer    string                 `json:"iss,omitempty"`
    Subject   string                 `json:"sub,omitempty"`
    Audience  Audience               `json:"aud,omitempty"`
    ExpiresAt *time.Time             `json:"exp,omitempty"`
    NotBefore *time.Time             `json:"nbf,omitempty"`
    IssuedAt  *time.Time             `json:"iat,omitempty"`
    JWTID     string                 `json:"jti,omitempty"`
    // OIDC profile fields
    Email     string                 `json:"email,omitempty"`
    Name      string                 `json:"name,omitempty"`
    Picture   string                 `json:"picture,omitempty"`
    Roles     []string               `json:"roles,omitempty"`
    Private   map[string]interface{} `json:"-"`
}
```

```go
// security/jwt/validator.go — claimsFromToken additions
if email, ok := mapClaims["email"].(string); ok {
    mapped.Email = email
}
if name, ok := mapClaims["name"].(string); ok {
    mapped.Name = name
}
if picture, ok := mapClaims["picture"].(string); ok {
    mapped.Picture = picture
}
// Roles: handle []interface{} (jwt lib default) and []string
switch r := mapClaims["roles"].(type) {
case []interface{}:
    for _, v := range r {
        if s, ok := v.(string); ok {
            mapped.Roles = append(mapped.Roles, s)
        }
    }
case []string:
    mapped.Roles = append([]string(nil), r...)
}
```

Private skip-list extended: `"email", "name", "picture", "roles"`.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `claimsFromToken` with OIDC fields present/absent | Table-driven in `validator_test.go`; mint tokens with `mustSignToken` |
| Unit | `Roles` with `[]interface{}` (jwt default) and `[]string` | Separate table rows |
| Unit | Fields absent → zero values; `Roles` absent → `nil` | Explicit nil/zero assertion |
| Unit | OIDC keys do NOT appear in `Private` | Assert `mapped.Private` does not contain extracted keys |

## Migration / Rollout

No migration required. New fields are additive. Callers using `Private["email"]` etc. continue to compile but will receive `nil` for those keys after this change — acceptable because typed fields are the replacement. No interface changes.

## Open Questions

- None
