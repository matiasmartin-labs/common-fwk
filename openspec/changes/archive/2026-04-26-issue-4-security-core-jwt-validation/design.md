# Design: Issue #4 Security Core JWT Validation

## Technical Approach

Implement an additive, framework-agnostic JWT validation core split into `security/claims`, `security/keys`, and `security/jwt`. The validator composes injected dependencies (key resolver, clock, parser options) and returns wrapped, assertable errors. `config` remains configuration-only; compatibility is handled by a thin mapping function from `config.JWTConfig` to validator options.

## Architecture Decisions

| Decision | Options | Tradeoff | Selected |
|---|---|---|---|
| Package boundaries | Single `security/jwt`; split `claims` + `keys` + `jwt` | Single package is simpler but mixes policy, claims, and key concerns; split adds files but keeps interfaces minimal and explicit | Split packages with narrow contracts |
| Key resolution | Inline key param; resolver interface | Inline key is easy but blocks `kid` selection and multi-key verification; resolver adds one abstraction but keeps deterministic tests | `keys.Resolver` interface keyed by `kid` |
| Time handling | `time.Now()` direct; injected clock | Direct time is simpler but non-deterministic; injected clock improves deterministic tests | `Now func() time.Time` option with default |
| Error model | String/opaque errors; sentinel+typed wrappers | Opaque errors are brittle; typed+sentinel errors require discipline but preserve `errors.Is/As` | Sentinel categories plus typed context wrappers |

## Data Flow

```text
Validate(token, opts)
  -> Parse JWT structure (claims+header)
  -> Method gate (alg allowlist)
  -> Resolve verification key (kid/default via resolver)
  -> Verify signature
  -> Claim checks (iss, aud, exp, nbf)
  -> Option checks (required issuer/audience policy)
  -> Return normalized claims
```

Validation order is intentional: reject malformed/disallowed method before expensive or ambiguous checks; classify failures by stage.

## File Changes

| File | Action | Description |
|---|---|---|
| `security/claims/doc.go` | Create | Package contract and scope. |
| `security/claims/claims.go` | Create | Standard claims model, audience normalization helpers. |
| `security/claims/claims_test.go` | Create | `aud` string/array normalization and optional claim behavior. |
| `security/keys/doc.go` | Create | Package contract for verification keys. |
| `security/keys/types.go` | Create | Key type and metadata (`kid`, method). |
| `security/keys/resolver.go` | Create | `Resolver` interface and deterministic in-memory resolver. |
| `security/keys/resolver_test.go` | Create | Present/missing key and default-key behavior. |
| `security/jwt/doc.go` | Create | Validator scope and non-goals. |
| `security/jwt/errors.go` | Create | Sentinel categories and typed validation error wrapper. |
| `security/jwt/options.go` | Create | Validator options (`Methods`, `Issuer`, `Audience`, `Now`, resolver). |
| `security/jwt/validator.go` | Create | Parse/validate orchestration and wrapping rules. |
| `security/jwt/validator_test.go` | Create | Deterministic validator contract tests (happy/failure paths). |
| `security/jwt/compat.go` | Create | Mapping helper from `config.JWTConfig` to jwt options. |
| `README.md` | Modify | Security core usage and compatibility notes. |

## Interfaces / Contracts

```go
// security/keys
type Key struct {
    ID     string
    Method string
    Verify any // concrete key material for jwt library
}

type Resolver interface {
    Resolve(ctx context.Context, kid string) (Key, error)
}

// security/jwt
type Options struct {
    Methods  []string
    Issuer   string
    Audience []string
    Now      func() time.Time
    Resolver keys.Resolver
}

type Validator interface {
    Validate(ctx context.Context, raw string) (claims.Claims, error)
}
```

Error taxonomy (`security/jwt/errors.go`):
- `ErrMalformedToken`
- `ErrInvalidSignature`
- `ErrInvalidIssuer`
- `ErrInvalidAudience`
- `ErrInvalidMethod`
- `ErrExpiredToken`
- `ErrNotYetValidToken`
- `ErrKeyResolution`

Wrapping rule: stage-specific failures are wrapped with `%w`; exported sentinels must remain assertable with `errors.Is`, and typed context (`ValidationError{Stage, Err}`) must remain assertable via `errors.As`.

## Testing Strategy

| Layer | What to Test | Approach |
|---|---|---|
| Unit (`security/claims`) | `aud` normalization and optional claims | Table tests for string/array/missing `aud`. |
| Unit (`security/keys`) | Resolver deterministic behavior | In-memory resolver tests for `kid` hit/miss and default fallback. |
| Unit (`security/jwt`) | Flow categories by stage | Fixed tokens + fixed clock + fake resolver; assert sentinels and typed wrappers. |
| Contract | `errors.Is/As` compatibility | Wrap returned errors and assert category/typed extraction still works. |

Determinism hooks are first-class: validator accepts injected `Now` and `Resolver`; tests avoid global time or network key fetch.

## Migration / Rollout

No breaking migration required. `config.ValidateConfig` semantics and `config.JWTConfig` fields (`Secret`, `Issuer`, `TTLMinutes`) stay unchanged.

Compatibility mapping contract:
- `JWTConfig.Secret` -> local symmetric verification key in resolver
- `JWTConfig.Issuer` -> `Options.Issuer`
- `JWTConfig.TTLMinutes` -> token-issuing concern (not validator gate), documented as compatible but not revalidated in `config`

Rollout is phased and additive: introduce security packages, add compatibility helper, then document usage.

## Open Questions

- [ ] Should initial method allowlist default to HS256-only or require explicit configuration?
- [ ] Should `compat.go` live in `security/jwt` or an adapter package if more runtimes appear?
