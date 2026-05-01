# Design: RS256 Keypair Security Bootstrap

## Technical Approach

Implement RS256 as an additive path over existing HS256 behavior. Keep `config.JWTConfig` backward-compatible (`HS256` default), add RS256-specific settings, and make validation conditional by algorithm. Reuse existing `security/jwt` method allowlist + resolver model by extending `security/jwt.FromConfigJWT` to produce either HS256 or RS256 `Options`. Add a small in-memory keypair utility in `security/keys` for deterministic bootstrap, then add an app convenience method that wires validator-from-config without replacing explicit `UseServerSecurity`.

## Architecture Decisions

### Decision: Single expanded JWT config model

| Option | Tradeoff | Decision |
|---|---|---|
| Split HS/RS config structs | Strong typing, but high migration churn across constructors/adapter/tests/docs | ❌ |
| Expand existing `JWTConfig` with algorithm + RS256 fields | Slightly richer struct, but additive and migration-safe | ✅ |

Rationale: Minimizes API churn and preserves existing call sites using `NewJWTConfig(secret, issuer, ttlMinutes)`.

### Decision: Deterministic in-memory RSA keypair API in `security/keys`

| Option | Tradeoff | Decision |
|---|---|---|
| Add provider/JWKS integration | More featureful, violates current non-goals and boundaries | ❌ |
| Keep deterministic local keypair generation/retrieval | Limited scope, aligns with core contracts and testability | ✅ |

Rationale: Keeps `security/*` provider-agnostic and enables robust unit tests for valid/invalid RS256 bootstrap.

### Decision: App bootstrap convenience as thin wrapper

| Option | Tradeoff | Decision |
|---|---|---|
| Replace `UseServerSecurity` with implicit global wiring | Simpler API, but boundary erosion and hidden state risk | ❌ |
| Add explicit convenience method delegating to `UseServerSecurity` | Slight extra method surface, but preserves explicit injection path | ✅ |

Rationale: Satisfies usability goal while keeping adapters dependent on core contracts and avoiding singleton key stores.

## Data Flow

1. Config is built (explicit constructors or `config/viper.Load`) and validated.
2. `security/jwt.FromConfigJWT(cfg.Security.Auth.JWT)` branches:
   - `HS256`: build static HMAC resolver from `secret`.
   - `RS256`: obtain resolver from configured RS256 key source (`generated` or provided PEM/public key).
3. App convenience method creates validator and calls `UseServerSecurity`.
4. HTTP middleware validates incoming tokens with algorithm-constrained options.

```text
config/viper.Load or constructors
           │
           ▼
      ValidateConfig
           │
           ▼
 security/jwt.FromConfigJWT
   ├─ HS256 ──> NewStaticResolver(HMAC)
   └─ RS256 ──> keys keypair/public resolver
           │
           ▼
  jwt.NewValidator(options) ──> app.UseServerSecurity
```

## File Changes

| File | Action | Description |
|---|---|---|
| `config/types.go` | Modify | Extend `JWTConfig` with algorithm + RS256 sub-config types. |
| `config/constructors.go` | Modify | Keep `NewJWTConfig` compatibility and add focused RS256 helpers/defaults. |
| `config/validate.go` | Modify | Conditional validation: secret required only for HS256; RS256 key-source constraints. |
| `config/validate_test.go` | Modify | Add HS256/RS256 valid/invalid matrix tests. |
| `config/viper/mapping.go` | Modify | Map canonical RS256 keys into expanded core JWT config. |
| `config/viper/loader.go` | Modify | Add env overrides + legacy alias handling for new JWT keys. |
| `config/viper/*_test.go` | Modify | Deterministic adapter tests for canonical/legacy/env/error cases. |
| `security/keys/keypair.go` | Create | RSA keypair generation/retrieval helpers and resolver bridge. |
| `security/keys/keypair_test.go` | Create | Deterministic behavior + error classification tests. |
| `security/jwt/compat.go` | Modify | Add HS256/RS256 branching for config-driven validator options. |
| `security/jwt/compat_test.go` | Create | Verify options/resolver wiring for both algorithms + invalid configs. |
| `app/application.go` | Modify | Add `UseServerSecurityFromConfig()` convenience wiring with wrapped errors. |
| `app/application_test.go` | Modify | Ordering/behavior tests for convenience method and failure modes. |
| `README.md`, `docs/home.md`, `docs/migration/auth-provider-ms-v0.1.0.md`, `docs/releases/v0.2.0-checklist.md` | Modify | Document mode semantics, bootstrap examples, migration/release impact. |

## Interfaces / Contracts

```go
// config/types.go
type JWTConfig struct {
    Algorithm  string // default: "HS256"
    Secret     string // required when Algorithm=HS256
    Issuer     string
    TTLMinutes int
    RS256      RS256Config
}

type RS256Config struct {
    KeySource    string // "generated" | "public-pem" | "private-pem"
    KeyID        string
    PublicKeyPEM string
    PrivateKeyPEM string
}

// security/keys/keypair.go
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, error)
func ResolverFromRS256Config(cfg config.RS256Config) (keys.Resolver, error)

// app/application.go
func (a *Application) UseServerSecurityFromConfig() (*Application, error)
```

## Testing Strategy

| Layer | What to Test | Approach |
|---|---|---|
| Unit | JWT config validation matrix | Table-driven tests for HS256/RS256 + invalid combinations. |
| Unit | Keypair helper behavior | Deterministic resolver output, nil/invalid PEM errors, 2048-bit default. |
| Unit | JWT compat mapping | Assert `Methods`, issuer, resolver type for HS256 and RS256. |
| Integration | Viper load + env aliasing | Load fixtures (kebab + camel legacy), assert canonical precedence and typed failures. |
| Integration | App convenience bootstrap | Verify success path and wrapped errors without replacing explicit API. |

## Migration / Rollout

No data migration required. Rollout is additive: existing HS256 configs continue unchanged; RS256 is opt-in via `security.auth.jwt.algorithm`. Documentation and release checklist updates ship in the same change.

## Open Questions

- [ ] Should `RS256.KeySource=generated` cache the generated keypair per app instance only, or per config bootstrap call? (Design assumes per bootstrap call, no global cache.)
