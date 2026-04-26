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

`./config/app.yaml` example:

```yaml
server:
  host: 127.0.0.1
  port: 8080
security:
  auth:
    jwt:
      secret: ${JWT_SECRET}
      issuer: common-fwk
      ttlMinutes: 15
    cookie:
      name: session
      domain: example.com
      secure: true
      httpOnly: true
      sameSite: Lax
    login:
      email: admin@example.com
    oauth2:
      providers:
        github:
          clientID: ${GITHUB_CLIENT_ID}
          clientSecret: ${GITHUB_CLIENT_SECRET}
          authURL: https://github.com/login/oauth/authorize
          tokenURL: https://github.com/login/oauth/access_token
          redirectURL: https://app.example.com/auth/github/callback
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
	compat := securityjwt.FromConfigJWT(jwtCfg)

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
