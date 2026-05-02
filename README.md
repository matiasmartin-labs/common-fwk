# common-fwk

Reusable framework primitives for config, security validation, and HTTP adapters.

## Quickstart

### Install

```bash
go get github.com/matiasmartin-labs/common-fwk@latest
```

### 1) Build validated core config (explicit API)

```go
package main

import (
	"errors"
	"fmt"

	"github.com/matiasmartin-labs/common-fwk/config"
)

func main() {
	cfg := config.NewConfig(
		config.NewServerConfig("127.0.0.1", 8080),
		config.NewSecurityConfig(
			config.NewAuthConfig(
				config.NewJWTConfig("secret", "common-fwk", 15),
				config.NewCookieConfig("session", "example.com", true, true, "Lax"),
				config.NewLoginConfig("  ADMIN@Example.COM  "),
				config.NewOAuth2Config(map[string]config.OAuth2ProviderConfig{
					"github": config.NewOAuth2ProviderConfig(
						"client-id",
						"client-secret",
						"https://github.com/login/oauth/authorize",
						"https://github.com/login/oauth/access_token",
						"https://app.example.com/auth/github/callback",
						[]string{"read:user", "user:email"},
					),
				}),
			),
		),
	)

	validated, err := config.ValidateConfig(cfg)
	if err != nil {
		if errors.Is(err, config.ErrInvalidConfig) {
			panic("invalid config")
		}
		panic(err)
	}

	fmt.Println(validated.Security.Auth.Login.Email)
	// Output: admin@example.com
}
```

### 2) Optional facade adapter (`config/viper`)

```go
package main

import (
	"errors"
	"fmt"

	"github.com/matiasmartin-labs/common-fwk/config"
	cfgviper "github.com/matiasmartin-labs/common-fwk/config/viper"
)

func main() {
	loaded, err := cfgviper.Load(cfgviper.Options{
		ConfigPath:  "./config/app.yaml",
		ConfigType:  "",           // inferred from extension when empty
		EnvPrefix:   "COMMON_FWK", // default when omitted
		EnvOverride: true,           // env values override file values
		ExpandEnv:   true,           // expands ${VAR} using a per-load env snapshot
	})
	if err != nil {
		var loadErr *cfgviper.LoadError
		var decodeErr *cfgviper.DecodeError
		var mapErr *cfgviper.MappingError
		var validateErr *cfgviper.ValidationError

		switch {
		case errors.As(err, &loadErr):
			fmt.Println("load error")
		case errors.As(err, &decodeErr):
			fmt.Println("decode error")
		case errors.As(err, &mapErr):
			fmt.Println("mapping error")
		case errors.As(err, &validateErr):
			fmt.Println("validation error")
		}

		if errors.Is(err, config.ErrInvalidConfig) {
			fmt.Println("core validation failed")
		}
		return
	}

	fmt.Println(loaded.Server.Host, loaded.Server.Port)
}
```

Behavior notes:
- `EnvOverride=false`: file values remain source of truth.
- `EnvOverride=true`: uses deterministic keys (example: `COMMON_FWK_SECURITY_AUTH_JWT_SECRET`).
- `ExpandEnv=false`: `${VAR}` placeholders are preserved.
- `ExpandEnv=true`: placeholders are expanded from a per-load env snapshot.
- File-based adapter keys are documented in canonical kebab-case.
- Legacy camelCase file keys remain temporarily compatible during migration.

Server runtime limits (core + adapter):
- File keys: `server.read-timeout`, `server.write-timeout`, `server.max-header-bytes`
- Defaults: `10s`, `10s`, `1048576` (1 MiB)
- Env override keys:
  - `COMMON_FWK_SERVER_READ_TIMEOUT`
  - `COMMON_FWK_SERVER_WRITE_TIMEOUT`
  - `COMMON_FWK_SERVER_MAX_HEADER_BYTES`

`./config/app.yaml` example:

```yaml
server:
  host: 127.0.0.1
  port: 8080
  read-timeout: 10s
  write-timeout: 10s
  max-header-bytes: 1048576
security:
  auth:
    jwt:
      secret: ${JWT_SECRET}
      issuer: common-fwk
      ttl-minutes: 15
    cookie:
      name: session
      domain: example.com
      secure: true
      http-only: true
      same-site: Lax
    login:
      email: admin@example.com
    oauth2:
      providers:
        github:
          client-id: ${GITHUB_CLIENT_ID}
          client-secret: ${GITHUB_CLIENT_SECRET}
          auth-url: https://github.com/login/oauth/authorize
          token-url: https://github.com/login/oauth/access_token
          redirect-url: https://app.example.com/auth/github/callback
          scopes: [read:user, user:email]
```

### 3) Security validator usage (core + compatibility facade)

```go
package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/matiasmartin-labs/common-fwk/config"
	securityjwt "github.com/matiasmartin-labs/common-fwk/security/jwt"
)

func main() {
	jwtCfg := config.NewJWTConfig("secret", "common-fwk", 15)
	compat, err := securityjwt.FromConfigJWT(jwtCfg)
	if err != nil {
		panic(err)
	}

	validator, err := securityjwt.NewValidator(compat.Options)
	if err != nil {
		panic(err)
	}

	validatedClaims, err := validator.Validate(context.Background(), "raw-token")
	if err != nil {
		if errors.Is(err, securityjwt.ErrInvalidMethod) {
			fmt.Println("invalid method")
		}
		if errors.Is(err, securityjwt.ErrExpiredToken) {
			fmt.Println("expired token")
		}

		var vErr *securityjwt.ValidationError
		if errors.As(err, &vErr) {
			fmt.Println("validation stage error")
		}
		return
	}

	_ = compat.TokenTTL // token issuing concern, not validator runtime enforcement
	fmt.Println(validatedClaims.Subject)
}
```

HS256 remains the default when `jwt.algorithm` is omitted.

RS256 config example (kebab-case adapter keys):

```yaml
security:
  auth:
    jwt:
      algorithm: RS256
      issuer: common-fwk
      ttl-minutes: 15
      rs256-key-source: private-pem
      rs256-key-id: auth-main
      rs256-private-key-pem: |
        -----BEGIN RSA PRIVATE KEY-----
        ...
        -----END RSA PRIVATE KEY-----
```

Supported RS256 key sources:
- `generated`: in-memory RSA keypair generated at bootstrap
- `public-pem`: resolver built from `rs256-public-key-pem`
- `private-pem`: resolver built from `rs256-private-key-pem`

### 4) Gin integration example (fluent option-style API)

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginfwk "github.com/matiasmartin-labs/common-fwk/http/gin"
	"github.com/matiasmartin-labs/common-fwk/security/keys"
	securityjwt "github.com/matiasmartin-labs/common-fwk/security/jwt"
)

func main() {
	validator, err := securityjwt.NewValidator(securityjwt.Options{
		Methods: []string{"HS256"},
		Issuer:  "common-fwk",
		Resolver: keys.NewStaticResolver(
			&keys.Key{Method: "HS256", Verify: []byte("secret")},
			nil,
		),
	})
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(ginfwk.NewAuthMiddleware(
		validator,
		ginfwk.WithHeaderName("Authorization"),
		ginfwk.WithCookieName("session"),
		ginfwk.WithContextKey("claims"),
		ginfwk.WithAuthEnabled(true),
	))

	r.GET("/me", func(c *gin.Context) {
		cl, ok := ginfwk.GetClaims(c, "claims")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "auth_token_missing"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"sub": cl.Subject})
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
```

### 5) Application bootstrap boundary (`app`)

`common-fwk/app` now includes an instance-based bootstrap API (no package-level singleton):

```go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matiasmartin-labs/common-fwk/app"
	"github.com/matiasmartin-labs/common-fwk/config"
	"github.com/matiasmartin-labs/common-fwk/security/keys"
	securityjwt "github.com/matiasmartin-labs/common-fwk/security/jwt"
)

func main() {
	cfg := config.NewConfig(
		config.NewServerConfig("127.0.0.1", 8080),
		config.NewSecurityConfig(config.NewAuthConfig(
			config.NewJWTConfig("secret", "common-fwk", 15),
			config.NewCookieConfig("session", "example.com", true, true, "Lax"),
			config.NewLoginConfig("admin@example.com"),
			config.NewOAuth2Config(nil),
		)),
	)

	validator, err := securityjwt.NewValidator(securityjwt.Options{
		Methods: []string{"HS256"},
		Issuer:  "common-fwk",
		Resolver: keys.NewStaticResolver(
			&keys.Key{Method: "HS256", Verify: []byte("secret")},
			nil,
		),
	})
	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApplication().
		UseConfig(cfg).
		UseServer().
		UseServerSecurity(validator)

	_ = a.RegisterGET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	_ = a.RegisterProtectedGET("/me", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Behavior notes:
- Register methods return typed sentinel errors when used out of order (for example, server/security not ready).
- `RegisterProtectedGET` uses `http/gin.NewAuthMiddleware` internally.
- `RunListener(net.Listener)` is available for test-friendly startup flows.
- `UseServerSecurityFromConfig()` is available as a thin convenience wrapper around config-driven JWT validator wiring (HS256/RS256).

Health/readiness preset API (explicit opt-in):

- `EnableHealthReadinessPresets(opts HealthReadinessOptions) error`

```go
package main

import (
	"errors"
	"log"

	"github.com/matiasmartin-labs/common-fwk/app"
	"github.com/matiasmartin-labs/common-fwk/config"
)

func main() {
	a := app.NewApplication().
		UseConfig(config.NewConfig(
			config.NewServerConfig("127.0.0.1", 8080),
			config.NewSecurityConfig(config.NewAuthConfig(
				config.NewJWTConfig("secret", "common-fwk", 15),
				config.NewCookieConfig("session", "example.com", true, true, "Lax"),
				config.NewLoginConfig("admin@example.com"),
				config.NewOAuth2Config(nil),
			)),
		)).
		UseServer()

	if err := a.EnableHealthReadinessPresets(app.HealthReadinessOptions{}); err != nil {
		if errors.Is(err, app.ErrServerNotReady) {
			log.Fatal("server bootstrap incomplete")
		}
		log.Fatal(err)
	}

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Custom paths + synchronous checks:

```go
package main

import (
	"errors"
	"log"

	"github.com/matiasmartin-labs/common-fwk/app"
	"github.com/matiasmartin-labs/common-fwk/config"
)

func main() {
	a := app.NewApplication().
		UseConfig(config.NewConfig(
			config.NewServerConfig("127.0.0.1", 8080),
			config.NewSecurityConfig(config.NewAuthConfig(
				config.NewJWTConfig("secret", "common-fwk", 15),
				config.NewCookieConfig("session", "example.com", true, true, "Lax"),
				config.NewLoginConfig("admin@example.com"),
				config.NewOAuth2Config(nil),
			)),
		)).
		UseServer()

	err := a.EnableHealthReadinessPresets(app.HealthReadinessOptions{
		HealthPath: "/livez",
		ReadyPath:  "/readyz/internal",
		Checks: []app.ReadinessCheck{
			func() error { return nil },
			func() error {
				if false {
					return errors.New("dependency not ready")
				}
				return nil
			},
		},
	})
	if err != nil {
		if errors.Is(err, app.ErrRouteConflict) {
			log.Fatal("preset route conflicts with existing GET route")
		}
		if errors.Is(err, app.ErrInvalidPresetOptions) {
			log.Fatal("invalid health/readiness options")
		}
		log.Fatal(err)
	}

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Readiness semantics:
- `/healthz` (or custom `HealthPath`) always returns `200` once presets are enabled.
- `/readyz` (or custom `ReadyPath`) returns `200` only when bootstrap invariants and all checks pass.
- Readiness returns `503` when invariants are unmet or any check fails.
- Presets are never auto-registered by `UseServer()`; default paths are not duplicated when custom paths are configured.

Health/readiness non-goals:
- No implicit preset registration during bootstrap.
- No provider-specific probing in framework internals; dependency readiness checks are caller-provided.

Read-only runtime accessors:

```go
fresh := app.NewApplication()

// Explicit non-init behavior.
_ = fresh.GetConfig()               // zero-value config snapshot
_ = fresh.GetSecurityValidator()    // nil
_ = fresh.IsSecurityReady()         // false

a := app.NewApplication().
	UseConfig(cfg).
	UseServer()

runtimeCfg := a.GetConfig()         // read-only snapshot
_ = runtimeCfg.Security.Auth.OAuth2.Providers // safe to inspect

// Mutating runtimeCfg does not mutate app internals.

a = a.UseServerSecurity(validator)
_ = a.GetSecurityValidator() != nil // true after security wiring
_ = a.IsSecurityReady()             // true
```

Accessor lifecycle + immutability contract:
- `GetConfig() config.Config`
- `GetSecurityValidator() security.Validator`
- `IsSecurityReady() bool`
- `GetRSAPrivateKey() *rsa.PrivateKey`
- `GetConfig()` is safe in all lifecycle stages and never triggers implicit bootstrap.
- Before security wiring, `GetSecurityValidator()` returns `nil` and `IsSecurityReady()` returns `false`.
- `GetConfig()` deep-copies mutable descendants (OAuth2 providers map and nested scopes slices) on each read.
- `GetRSAPrivateKey()` returns a non-nil key only when `UseServerSecurityFromConfig()` was called with RS256 algorithm and `Generated` or `PrivatePEM` key source. Returns `nil` for `PublicPEM`, direct `UseServerSecurity(v)` wiring, or when no security is wired. Never panics.

### 6) Named logger registry (`GetLogger`) and logging config

`app.Application` exposes deterministic named logger access:

- `GetLogger(name string) (logging.Logger, error)`

Lifecycle contract:
- Calling `GetLogger("auth")` before `UseConfig(...)` returns `ErrLoggingNotReady`.
- Calling `GetLogger("")` returns `ErrLoggerNameRequired`.
- Repeated `GetLogger("auth")` calls on the same `Application` return the same logger instance.
- Same logger name across different `Application` instances is isolated (different instances).

Core logging config keys:

- `logging.enabled` (default `true`)
- `logging.level` (default `info`) — accepted: `debug|info|warn|error`
- `logging.format` (default `json`) — accepted: `json|text`
- `logging.loggers.<name>.enabled` (optional override)
- `logging.loggers.<name>.level` (optional override)

Deterministic precedence matrix:

| Setting | Effective value |
|---|---|
| `enabled` | per-logger override when set, otherwise root `logging.enabled` |
| `level` | per-logger override when set, otherwise root `logging.level` |
| `format` | root `logging.format` only (`json|text`) |

`./config/app.yaml` logging example:

```yaml
logging:
  enabled: true
  level: info
  format: json
  loggers:
    auth:
      enabled: true
      level: debug
    billing:
      enabled: false
```

Environment override keys for logging (when `EnvOverride=true`):

- `COMMON_FWK_LOGGING_ENABLED`
- `COMMON_FWK_LOGGING_LEVEL`
- `COMMON_FWK_LOGGING_FORMAT`
- `COMMON_FWK_LOGGING_LOGGERS_<NAME>_ENABLED`
- `COMMON_FWK_LOGGING_LOGGERS_<NAME>_LEVEL`

Examples:
- `COMMON_FWK_LOGGING_LOGGERS_AUTH_LEVEL=debug`
- `COMMON_FWK_LOGGING_LOGGERS_BILLING_ENABLED=false`

Required output contract for accepted records:
- `logger`
- `ts`
- `level`
- `msg`

JSON output example:

```json
{"ts":"2026-05-01T19:40:00Z","level":"INFO","msg":"request accepted","logger":"auth"}
```

Text output example:

```text
ts=2026-05-01T19:40:00Z level=INFO msg="request accepted" logger=auth
```

Loki guidance (collector-first):
- Keep application logs framework-structured and sink-agnostic.
- Prefer Promtail/OpenTelemetry collector pipelines to ship logs to Loki.
- Preserve structured keys (`logger`, `ts`, `level`, `msg`) end-to-end for queryability.

---

## Layered architecture overview

`common-fwk` is organized in layers to keep domain logic reusable and adapter-agnostic:

1. **Core domain layer**
   - `config`: typed config model + constructors + validation.
   - `security/claims`, `security/keys`, `security/jwt`: token validation model/contracts/engine.
2. **Adapter layer**
   - `config/viper`: optional config loader adapter into core `config.Config`.
   - `http/gin`: optional HTTP middleware adapter that depends on `security.Validator` interface.
3. **App boundary**
- `app`: application bootstrap boundary.
  - `Application` fluent bootstrap: `UseConfig`, `UseServer`, `UseServerSecurity`.
  - Route registration: `RegisterGET`, `RegisterPOST`, `RegisterProtectedGET`.
  - Runtime startup: `Run`, `RunListener`.

Dependency direction is always inward:
- Adapters depend on core contracts.
- Core packages do **not** depend on framework-specific adapters.

## Package responsibilities and boundaries

- `config`
  - owns the canonical configuration model and validation rules.
  - no global mutable state, no adapter runtime coupling.

- `config/viper`
  - optional loader facade for file + env decoding.
  - maps adapter-local raw structs to core types, then validates via `config.ValidateConfig`.

- `security/claims`
  - claim model + audience normalization helpers.

- `security/keys`
  - deterministic key resolution contracts (`kid` + default fallback), no network fetch.

- `security/jwt`
  - framework-agnostic JWT validator options, typed errors, validation flow.

- `http/gin`
  - thin adapter that extracts tokens, calls `security.Validator`, and injects claims into `gin.Context`.

## Non-goals

- No app-global singleton coupling.
- No framework lock-in inside core packages.
- No remote provider/JWKS fetcher coupling in `security/*`.
- **Google provider implementation is intentionally outside this framework**.
  - This repo provides generic OAuth2 config fields.
  - Concrete Google OAuth flows, APIs, and provider-specific behavior belong to the consuming application/service layer.

## Documentation

Full documentation is available under [`docs/`](docs/index.md) and published to GitHub Pages.

| Section | Path | Description |
|---|---|---|
| Home | `docs/index.md` | Landing page and quick reference |
| Getting Started | `docs/getting-started/` | Installation and quickstart |
| Architecture | `docs/architecture/` | Canonical specs for each package |
| Releases | `docs/releases/` | Release notes per version (v0.1.0 → v0.7.0) |
| Migration Guides | `docs/migration/` | Consumer migration guides |
| Decisions | `docs/decisions/` | Architecture decision records (ADRs) |
| Contributing | `docs/contributing/` | SDD workflow and PR conventions |

## Release and migration docs

- Release workflow label policy:
  - Use exactly one of `release-type:patch`, `release-type:minor`, or `release-type:major` on release PRs.
  - Use `release:skip` to explicitly skip release preview/publication for non-release changes (docs/chore-only PRs).
- Release notes: `docs/releases/` (v0.1.0 through v0.7.0)
- Migration guide for `auth-provider-ms`: `docs/migration/auth-provider-ms-v0.1.0.md`
- Documentation structure migration note: `docs/migration/openspec-to-docs.md`
