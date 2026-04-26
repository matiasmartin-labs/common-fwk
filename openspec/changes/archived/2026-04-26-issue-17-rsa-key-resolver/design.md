# Design: RSA Key Resolver for RS256 JWT Validation

## Technical Approach

Add RSA-focused resolver constructors in `security/keys` that reuse the existing `NewStaticResolver` and `Key{Method, Verify}` contract. `jwt.Validator` already forwards `Key.Verify` to `golang-jwt/jwt/v5`, which dispatches verification by concrete key type, so no validator/runtime contract changes are required.

Implementation is intentionally additive: one new file for constructors plus RS256 coverage in existing validator tests. HS256 paths remain unchanged.

## Architecture Decisions

| Option | Tradeoff | Decision |
|---|---|---|
| Add `PublicKey` field to `keys.Key` and update validator selection logic | More explicit type, but duplicates existing `Verify any` responsibility and increases risk in core validation path | Rejected |
| Keep `Key.Verify any` and add RSA constructor wrappers | Minimal, backward compatible, no validator changes | **Chosen** |

| Option | Tradeoff | Decision |
|---|---|---|
| Only `NewRSAResolver(*rsa.PrivateKey, keyID string)` | Small API, but awkward for verify-only services that only have public key | Rejected |
| Add both private-key and public-key constructors | Slightly larger surface, clearer caller intent and broader usability | **Chosen** |

| Option | Tradeoff | Decision |
|---|---|---|
| Validate nil keys via returned error | Best UX, but conflicts with required constructor signature returning only `Resolver` | Rejected |
| Keep constructor total (non-panicking), propagate invalid-key failures at validation time | Aligns with existing constructor style (`NewStaticResolver`), avoids panic and signature changes | **Chosen** |

## Data Flow

`RS256 token` → `Validator.Validate` → parse unverified header (`alg`, `kid`) → `Resolver.Resolve(ctx, kid)` → `keys.Key{Method:"RS256", Verify:*rsa.PublicKey or *rsa.PrivateKey}` → jwt parser keyfunc returns `Key.Verify` → `golang-jwt/v5` verifies signature by key type → claims checks (`iss`, `aud`, `exp`, `nbf`) → result.

```
Caller config
  └─ keys.NewRSAPublicKeyResolver(pub, keyID)
        └─ NewStaticResolver(default RS256 key)

Validator.Validate(raw)
  ├─ Resolve(kid) -> Key
  ├─ keyfunc -> Key.Verify
  └─ golang-jwt RS256 verify + existing claim validation
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `security/keys/rsa.go` | Create | Add `NewRSAResolver` and `NewRSAPublicKeyResolver` wrappers around `NewStaticResolver` using method `RS256`. |
| `security/jwt/validator_test.go` | Modify | Add RS256 scenarios: valid token, invalid signature, expired token; keep existing HS256 scenarios unchanged. |

No planned changes: `security/keys/types.go`, `security/jwt/validator.go`, `security/jwt/options.go`.

## Interfaces / Contracts

```go
// security/keys/rsa.go
func NewRSAResolver(privateKey *rsa.PrivateKey, keyID string) Resolver
func NewRSAPublicKeyResolver(publicKey *rsa.PublicKey, keyID string) Resolver
```

Contract notes:
- Both constructors return deterministic static resolvers with default key only.
- `Key.Method` is set to `"RS256"`.
- `NewRSAPublicKeyResolver` sets `Key.Verify` to `*rsa.PublicKey`.
- `NewRSAResolver` sets `Key.Verify` to RSA material compatible with RS256 verification and avoids panic on nil input.
- Callers MUST set `jwt.Options.Methods` to include `"RS256"` (default remains `HS256`).

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | RS256 success path | Generate RSA keypair in test, sign RS256 token, validate with `Methods:["RS256"]` and `NewRSAPublicKeyResolver`. |
| Unit | RS256 invalid signature | Sign with different private key than resolver key; assert `errors.Is(err, ErrInvalidSignature)`. |
| Unit | RS256 expiration handling | RS256 token with `exp` in the past; assert `errors.Is(err, ErrExpiredToken)`. |
| Regression | Existing HS256 behavior | Keep existing test table and assertions intact; run full package tests. |

## Migration / Rollout

No migration required. Rollout is additive: consumers opt in by constructing RSA resolvers and allowing RS256 in validator options.

## Open Questions

- [ ] None.
