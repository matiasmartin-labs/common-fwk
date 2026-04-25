# Design: Issue #2 Config Core

## Technical Approach

Implement a small, typed `config` core package with explicit constructors, deterministic validation, and a stable error model, while keeping adapter concerns (`config/viper`) out of scope. The implementation starts by narrowing bootstrap guards so `config/` can evolve beyond `doc.go` without weakening bootstrap-only protections for other packages.

## Architecture Decisions

| Decision | Options | Tradeoffs | Selected |
|---|---|---|---|
| Core model boundaries | Flat root struct vs nested domain structs | Flat is shorter but blurs ownership; nested keeps server/security/auth boundaries clear | Nested: `Config{Server, Security}` and `SecurityConfig{Auth}` |
| Constructor style | Public literals only vs `New*` constructors + defaults | Literals are flexible but spread defaults; constructors centralize invariants and defaults | `NewConfig` + focused `New*Config` helpers returning values |
| Error model | String-only errors vs sentinels + typed path errors | String-only is brittle in tests; typed model is slightly more code but assertable | `ErrXxx` sentinels + `ValidationError{Path, Err}` wrapping |
| Validation composition | One monolithic validator vs subtree validators | Monolith is simple initially but grows brittle; subtree validators stay readable and testable | Root orchestration calling per-subtree validators |
| Login normalization timing | Normalize only in constructors vs during validation | Constructor-only can miss literal-built configs; validation-time ensures consistency | Normalize in validation entrypoint before rule checks |

## Data Flow

`NewConfig(...)` (or literal assembly) → `ValidateConfig(cfg)` → normalize login email → run subtree validators (`server`, `jwt`, `cookie`, `login`, `oauth2`) → return normalized `Config` + wrapped error.

```
Caller
  │
  ├─ construct (New* / literals)
  ▼
Config value
  │
  ├─ ValidateConfig
  │    ├─ normalizeLoginEmail
  │    ├─ validateServer
  │    ├─ validateJWT
  │    ├─ validateCookie
  │    ├─ validateLogin
  │    └─ validateOAuth2
  ▼
normalized Config, error (nil or wrapped/typed)
```

## File Changes

| File | Action | Description |
|---|---|---|
| `bootstrap_guard_test.go` | Modify | Remove `config` from doc-only package list; keep `config/viper` and other bootstrap packages guarded. Add assertion comment to preserve intent. |
| `config/types.go` | Create | Define `Config`, `ServerConfig`, `SecurityConfig`, `AuthConfig`, `JWTConfig`, `CookieConfig`, `LoginConfig`, `OAuth2Config`, and provider-client generic type(s). |
| `config/constructors.go` | Create | `NewConfig` and focused `New*` helpers with useful zero/default behavior and no globals. |
| `config/errors.go` | Create | Sentinel errors (`ErrInvalidConfig`, `ErrRequired`, `ErrOutOfRange`, `ErrInvalidEmail`, etc.) and `ValidationError` wrapper implementing `Unwrap`. |
| `config/validate.go` | Create | `ValidateConfig(cfg Config) (Config, error)` plus internal composable validators per subtree. |
| `config/constructors_test.go` | Create | Deterministic constructor/default tests and zero-value behavior checks. |
| `config/validate_test.go` | Create | Table-driven valid/invalid cases, normalization checks, and `errors.Is` / `errors.As` assertions. |
| `config/errors_test.go` | Create | Error wrapping and path metadata behavior tests. |
| `config/doc.go` | Modify | Package contract + no-global/no-adapter statement. |
| `README.md` | Modify | Minimal usage snippet for construction + validation.

## Interfaces / Contracts

```go
package config

type Config struct { Server ServerConfig; Security SecurityConfig }
type ServerConfig struct { Host string; Port int }
type SecurityConfig struct { Auth AuthConfig }
type AuthConfig struct { JWT JWTConfig; Cookie CookieConfig; Login LoginConfig; OAuth2 OAuth2Config }
type JWTConfig struct { Secret string; Issuer string; TTLMinutes int }
type CookieConfig struct { Name string; Domain string; Secure bool; HTTPOnly bool; SameSite string }
type LoginConfig struct { Email string }
type OAuth2Config struct { Providers map[string]OAuth2ProviderConfig }
type OAuth2ProviderConfig struct { ClientID string; ClientSecret string; AuthURL string; TokenURL string; RedirectURL string; Scopes []string }

func NewConfig(server ServerConfig, security SecurityConfig) Config
func ValidateConfig(cfg Config) (Config, error)

var (
    ErrInvalidConfig error
    ErrRequired error
    ErrOutOfRange error
    ErrInvalidEmail error
)

type ValidationError struct { Path string; Err error }
```

## Testing Strategy

| Layer | What to Test | Approach |
|---|---|---|
| Unit | Constructors/defaults and zero-value usefulness | Table-driven tests; compare structs with exact expected values. |
| Unit | Validation composition and error taxonomy | Subtests per subtree; deterministic fixtures; assert `errors.Is` sentinels and `errors.As(*ValidationError)` path. |
| Unit | Login normalization | Inputs with spaces/uppercase; assert returned config stores trimmed lowercase email. |
| Integration-lite | Bootstrap guard compatibility | Update existing guard test to ensure `config/` growth is allowed while other bootstrap packages still restricted. |

Determinism safeguards: no time/env/file dependencies, no global mutable state, stable map assertions via sorted provider keys before comparison.

## Migration / Rollout

No migration required. Rollout is additive in `config/` plus a targeted bootstrap guard adjustment.

## Open Questions

- [ ] Should `SameSite` be a string enum alias now or deferred to follow-up hardening?
