# Design: RSA Public Key and Key ID Accessors on app.Application

## Technical Approach

Mirror issue-52's `GetRSAPrivateKey()` pattern exactly. Extend `CompatOptions` with two new fields, populate them in `resolveRS256` for all three RS256 key sources, capture them in `Application` during `UseServerSecurityFromConfig`, and expose them as O(1) nil-safe accessors.

## Architecture Decisions

| Option | Tradeoff | Decision |
|--------|----------|----------|
| Add fields to `CompatOptions` | Minimal surface; consistent with `RSAPrivateKey` already there | ✅ Chosen |
| Return public key from resolver | Requires interface change; breaks provider boundary | ❌ Rejected |
| Recompute public key on each call | Wastes CPU; breaks nil semantics for HS256/unwired | ❌ Rejected |

## Data Flow

```
FromConfigJWT (RS256 branch)
  └─→ resolveRS256(cfg)
        ├─ Generated:   priv → &priv.PublicKey + cfg.KeyID
        ├─ PrivatePEM:  priv → &priv.PublicKey + cfg.KeyID
        └─ PublicPEM:   parsedPub + cfg.KeyID
  └─→ CompatOptions{ RSAPublicKey, RSAKeyID, RSAPrivateKey, ... }

UseServerSecurityFromConfig
  └─→ a.rsaPublicKey = compat.RSAPublicKey
  └─→ a.rsaKeyID     = compat.RSAKeyID

GetRSAPublicKey() → a.rsaPublicKey  (nil if HS256 or unwired)
GetRSAKeyID()     → a.rsaKeyID      ("" if HS256 or unwired)
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `security/jwt/compat.go` | Modify | Add `RSAPublicKey *rsa.PublicKey` + `RSAKeyID string` to `CompatOptions`; populate both in `resolveRS256` for all three branches |
| `app/application.go` | Modify | Add `rsaPublicKey *rsa.PublicKey` + `rsaKeyID string` fields; capture from `compat`; add two accessor methods |
| `app/application_test.go` | Modify | Add `TestGetRSAPublicKey` + `TestGetRSAKeyID` table-driven tests covering all 5 scenarios each |
| `app/doc.go` | Modify | Add `GetRSAPublicKey` + `GetRSAKeyID` signatures |
| `README.md` | Modify | Add accessor rows to accessor table |
| `docs/home.md` | Modify | Add accessor rows to accessor table |

## Interfaces / Contracts

```go
// CompatOptions — additions only
type CompatOptions struct {
    // ... existing fields ...
    RSAPublicKey *rsa.PublicKey // non-nil for all RS256 key sources
    RSAKeyID     string         // non-empty for all RS256 key sources
}

// Application — additions only
func (a *Application) GetRSAPublicKey() *rsa.PublicKey { return a.rsaPublicKey }
func (a *Application) GetRSAKeyID() string             { return a.rsaKeyID }
```

`resolveRS256` changes:

```go
// Generated & PrivatePEM branches: append to return
RSAPublicKey: &priv.PublicKey,
RSAKeyID:     cfg.KeyID,

// PublicPEM branch: set after parsing pub
RSAPublicKey: parsedPub,
RSAKeyID:     cfg.KeyID,
```

Note: `resolveRS256` currently returns `(*rsa.PrivateKey, keys.Resolver, error)`. To also return the public key without changing the signature, the caller (`FromConfigJWT`) will derive `RSAPublicKey` from `CompatOptions.RSAPrivateKey` for Generated/PrivatePEM and from a new return value for PublicPEM. Simplest: change `resolveRS256` to return `(*rsa.PrivateKey, *rsa.PublicKey, keys.Resolver, error)` — internal function, no external API break.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `GetRSAPublicKey` — all 5 scenarios | Table-driven in `application_test.go` |
| Unit | `GetRSAKeyID` — all 3 scenarios | Table-driven in `application_test.go` |
| Unit | `CompatOptions` field population | Existing `compat_test.go` extended |

## Migration / Rollout

No migration required. Purely additive — new fields on existing structs, new methods on existing type. Zero-value semantics (`nil`, `""`) are safe for unwired Applications.

## Open Questions

- None.
