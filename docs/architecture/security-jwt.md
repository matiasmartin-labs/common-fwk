---
title: JWT Security
parent: Architecture
nav_order: 4
---

# JWT Security (`security/jwt`, `security/claims`, `security/keys`)

## Packages

| Package | Import | Purpose |
|---|---|---|
| `security/jwt` | `.../security/jwt` | JWT validator |
| `security/claims` | `.../security/claims` | Claims model |
| `security/keys` | `.../security/keys` | Key resolver contracts and constructors |

## Claims Model (`security/claims`)

`Claims` exposes typed fields plus a private map for non-standard claims:

| Field | Type | Source JWT key |
|---|---|---|
| `Email` | `string` | `email` |
| `Name` | `string` | `name` |
| `Picture` | `string` | `picture` |
| `Roles` | `[]string` | `roles` (bare key only) |
| `Private` | `map[string]any` | All non-standard keys |

- `aud` accepts string or array — normalized consistently.
- Missing optional fields default to zero values without failing parsing.
- Namespaced keys (e.g. `https://example.com/roles`) land in `Private`, not in `Roles`.

## Validator (`security/jwt`)

```go
validator, err := jwt.NewValidator(jwtConfig, resolver, jwt.WithTimeSource(fixedTime))
claims, err := validator.Validate(tokenString)
```

### Policy checks

- Signature verification.
- Issuer (`iss`) must match config.
- Audience (`aud`) must match config.
- Algorithm (`alg`) must be in allowlist.
- `exp` and `nbf` evaluated using injected time source.

### Error categories

| Category | Sentinel |
|---|---|
| Malformed token | `jwt.ErrMalformed` |
| Invalid signature | `jwt.ErrInvalidSignature` |
| Invalid issuer | `jwt.ErrInvalidIssuer` |
| Invalid audience | `jwt.ErrInvalidAudience` |
| Invalid method | `jwt.ErrInvalidMethod` |
| Expired token | `jwt.ErrExpired` |
| Not yet valid | `jwt.ErrNotYetValid` |
| Key resolution failure | `jwt.ErrKeyResolution` |

All errors are assertable via `errors.Is`/`errors.As`.

## Key Resolvers (`security/keys`)

| Constructor | Usage |
|---|---|
| `keys.NewStaticResolver(secret []byte)` | HS256 shared secret |
| `keys.NewRSAResolver(privateKey, keyID)` | RS256 private key (signs and verifies) |
| `keys.NewRSAPublicKeyResolver(publicKey, keyID)` | RS256 public key (verify only) |

Resolvers are deterministic and in-memory — no network I/O.

Key selection is by `kid` header field, falling back to default if not present.

## HS256 Example

```go
resolver := keys.NewStaticResolver([]byte("my-secret"))
validator, err := jwt.NewValidator(cfg.Security.Auth.JWT, resolver)
```

## RS256 Example

```go
// From config (recommended)
// UseServerSecurityFromConfig() handles this automatically

// Manual
resolver := keys.NewRSAPublicKeyResolver(pubKey, "my-key-id")
opts := jwt.ValidatorOptions{
    Issuer:  "my-service",
    Methods: []string{"RS256"},
}
validator, err := jwt.NewValidatorWithOptions(opts, resolver)
```

## Boundaries

- No Gin dependency.
- No app globals.
- No JWKS/OAuth provider adapters.
