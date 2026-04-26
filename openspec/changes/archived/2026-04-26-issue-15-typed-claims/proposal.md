# Proposal: Typed OIDC Claims on `claims.Claims`

## Intent

`claims.Claims` currently exposes only 7 standard JWT fields; OIDC profile fields (`email`, `name`, `picture`, `roles`) fall into the untyped `Private map[string]interface{}`. Consumers must assert types manually, silently failing on mis-keyed lookups. Adding these as first-class typed fields eliminates runtime casts, improves discoverability, and aligns with OIDC profile scope conventions.

## Scope

### In Scope
- Add `Email string`, `Name string`, `Picture string`, `Roles []string` to `claims.Claims`
- Populate new fields in `jwt/validator.go` → `claimsFromToken()` from standard OIDC claim keys
- Unit tests for new fields (valid, missing, wrong type graceful handling)
- Update delta spec for `security-core-jwt-validation`

### Out of Scope
- Namespaced/issuer-specific `roles` variants (stay in `Private`)
- Generics-based accessor pattern (over-engineering, deferred)
- Changes to `gin-auth-middleware` or `http/gin/context.go`
- `jwks` / key provider changes

## Capabilities

### New Capabilities
- None

### Modified Capabilities
- `security-core-jwt-validation`: Claims model gains OIDC profile fields; `claimsFromToken` mapping behavior extends to new keys

## Approach

Extend the `Claims` struct with four typed fields. In `claimsFromToken`, after mapping registered claims, extract `email`, `name`, `picture`, and `roles` from the raw token map using type-safe assertions (zero-value on miss — no error). All unknown claims continue flowing to `Private`. No interface changes; fully backward compatible.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `security/claims/claims.go` | Modified | Add `Email`, `Name`, `Picture`, `Roles` fields |
| `security/jwt/validator.go` | Modified | Populate new fields in `claimsFromToken` |
| `security/claims/claims_test.go` | Modified | Tests for typed field access |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| `roles` key name varies by issuer | Med | Only map bare `roles`; others stay in `Private` |
| Breaking consumer code that reads these from `Private` | Low | Additive struct fields; Private still populated for unknown keys |

## Rollback Plan

Revert `claims.go` and `jwt/validator.go` to pre-change state. Both files are isolated; no downstream interface changes required.

## Dependencies

- None

## Success Criteria

- [ ] `Claims.Email`, `Claims.Name`, `Claims.Picture`, `Claims.Roles` accessible without type assertion
- [ ] `claimsFromToken` populates new fields from standard OIDC keys
- [ ] Unknown/non-standard claims still land in `Private`
- [ ] All existing tests pass unchanged
- [ ] New table-driven tests cover: field present, field absent (zero value), wrong type (zero value)
