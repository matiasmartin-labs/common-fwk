# Design: RSA Private Key Accessor

## Technical Approach

Surface `*rsa.PrivateKey` from key derivation through `CompatOptions` into `Application`,
exposing it via a nil-safe `GetRSAPrivateKey()` accessor. No changes to resolver internals,
validator chain, or key-parsing logic.

## Architecture Decisions

| Topic | Choice | Rejected | Rationale |
|-------|--------|----------|-----------|
| Where to carry the key | `CompatOptions.RSAPrivateKey *rsa.PrivateKey` | Return as 3rd value from `FromConfigJWT` | Keeps the function signature stable; CompatOptions already bundles issuing concerns (TokenTTL) |
| Application field | `rsaPrivateKey *rsa.PrivateKey` (unexported) | Exported field | Consistent with `validator`, `cfg` — access only via accessor |
| Accessor signature | `GetRSAPrivateKey() *rsa.PrivateKey` | `GetRSAPrivateKey() (*rsa.PrivateKey, error)` | Key absence is not an error; nil return is the idiomatic zero value |
| Populate in compat | Assign in `FromConfigJWT` RS256 case only | Assign in `ResolverFromRS256Config` | compat.go already owns the RS256 branch; keeps keys package unchanged |

## Data Flow

```
config.JWTConfig (RS256)
    │
    ▼
FromConfigJWT()  ──── keys.ResolverFromRS256Config(cfg.RS256)
    │                         │
    │                    [Generated]   priv ← GenerateRSAKeyPair()
    │                    [PrivatePEM]  priv ← parsePrivateKeyPEM()
    │                    [PublicPEM]   priv = nil (public-only source)
    │
    ▼
CompatOptions {
    Options       Options
    TokenTTL      time.Duration
    RSAPrivateKey *rsa.PrivateKey   ← new field (nil for HS256 / PublicPEM)
}
    │
    ▼
UseServerSecurityFromConfig()
    │  compat, err := FromConfigJWT(...)
    │  a.rsaPrivateKey = compat.RSAPrivateKey
    ▼
Application.rsaPrivateKey *rsa.PrivateKey
    │
    ▼
GetRSAPrivateKey() *rsa.PrivateKey   ← nil when not set
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `security/jwt/compat.go` | Modify | Add `RSAPrivateKey *rsa.PrivateKey` to `CompatOptions`; populate in RS256 `Generated` and `PrivatePEM` branches |
| `app/application.go` | Modify | Add `rsaPrivateKey *rsa.PrivateKey` field; capture from compat in `UseServerSecurityFromConfig`; add `GetRSAPrivateKey()` accessor |
| `app/doc.go` | Modify | Add `GetRSAPrivateKey()` to the runtime inspection helpers list and lifecycle contract |
| `app/application_test.go` | Modify | Add test cases (see Testing Strategy) |
| `README.md` | Modify | Document `GetRSAPrivateKey()` in accessors section |
| `docs/home.md` | Modify | Mirror README accessor entry for doc-sync test compliance |

## Interfaces / Contracts

```go
// security/jwt/compat.go

type CompatOptions struct {
    Options      Options
    TokenTTL     time.Duration
    RSAPrivateKey *rsa.PrivateKey // non-nil only for RS256 Generated/PrivatePEM sources
}
```

```go
// app/application.go

// GetRSAPrivateKey returns the RSA private key derived during security wiring.
//
// Returns nil when security was not wired, algorithm is not RS256, or the
// configured key source does not provide a private key (e.g. PublicPEM).
func (a *Application) GetRSAPrivateKey() *rsa.PrivateKey {
    return a.rsaPrivateKey
}
```

Population in `UseServerSecurityFromConfig` (after existing validator wiring):

```go
a.rsaPrivateKey = compat.RSAPrivateKey
```

Population in `FromConfigJWT` RS256 branch — requires refactoring the single
`ResolverFromRS256Config` call to extract the private key before building the resolver.
Since `ResolverFromRS256Config` is unchanged, the key must be derived inline for the two
private-key sources, **or** a helper is introduced. Preferred approach: extract key
derivation for `Generated` and `PrivatePEM` cases inline within `FromConfigJWT` so the
private key is available before constructing the resolver via `NewRSAResolver`:

```go
case config.JWTAlgorithmRS256:
    priv, resolver, err := resolveRS256(cfg.RS256)
    // priv is nil for PublicPEM source
    return CompatOptions{..., RSAPrivateKey: priv}, nil
```

Where `resolveRS256` is a private helper in `compat.go` that calls into `keys` package
functions directly and returns `(*rsa.PrivateKey, keys.Resolver, error)`.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `FromConfigJWT` RS256 Generated → `RSAPrivateKey` non-nil | Table test in `security/jwt/compat_test.go` |
| Unit | `FromConfigJWT` RS256 PrivatePEM → `RSAPrivateKey` non-nil | Same table |
| Unit | `FromConfigJWT` RS256 PublicPEM → `RSAPrivateKey` nil | Same table |
| Unit | `FromConfigJWT` HS256 → `RSAPrivateKey` nil | Same table |
| Integration | `UseServerSecurityFromConfig` RS256 Generated → `GetRSAPrivateKey()` non-nil | `app/application_test.go` |
| Integration | `UseServerSecurityFromConfig` HS256 → `GetRSAPrivateKey()` nil | `app/application_test.go` |
| Unit | `GetRSAPrivateKey()` before any security wiring → nil | `app/application_test.go` |

## Migration / Rollout

No migration required. `RSAPrivateKey` defaults to nil; all existing callers of
`CompatOptions` and `Application` are unaffected.

## Open Questions

- None. Design is fully determined by existing patterns.
