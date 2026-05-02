---
name: common-fwk-usage
description: >
  Create or update integration documentation and usage examples for common-fwk.
  Trigger: When writing README quickstarts, architecture overviews, package boundaries,
  or integration examples for config/security/http/app adapters.
license: MIT
metadata:
  author: Matias Martin
  version: "2.0"
---

## When to Use

- Updating README usage docs for `common-fwk`
- Writing quickstart examples that users should follow end-to-end
- Documenting architecture layers and dependency boundaries
- Documenting package responsibilities and explicit non-goals
- Bootstrapping an application using the `app` package lifecycle

## Critical Patterns

- Keep examples buildable: include `package main`, imports, and `main()` when snippet scope is full-flow.
- Prefer explicit core API first (`config`, `security/jwt`) then optional adapters (`config/viper`, `http/gin`, `app`).
- Keep boundary clarity explicit: adapters depend on core contracts, never the other way around.
- Preserve non-goals in docs: no app-global singletons, no framework lock-in in core, no remote provider coupling in `security/*`.
- State provider ownership clearly: Google OAuth provider logic stays in consuming apps, outside this framework.
- `NewConfig` logging parameter is variadic — it MAY be omitted.

## Documentation Structure

All current packages:

1. **`app`** — Application bootstrap: wires config, server, security, and health endpoints.
2. **`config`** — Core config constructors: server, security, auth, JWT, RS256, cookie, login, OAuth2, logging.
3. **`config/viper`** — File + env-based config loader using Viper.
4. **`errors`** — Shared error types used across the framework.
5. **`http/gin`** — Gin auth middleware with functional options.
6. **`logging`** — Logger and Registry interfaces; `logging/slog` adapter.
7. **`security`** — Core security validator interface.
8. **`security/claims`** — `claims.Claims` interface extracted from validated JWT tokens.
9. **`security/jwt`** — JWT validator supporting HS256 and RS256.
10. **`security/keys`** — Key resolvers: static (HMAC), RSA private, RSA public.

## Code Examples

### Example 1 — End-to-end app quickstart

```go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/matiasmartin-labs/common-fwk/app"
	"github.com/matiasmartin-labs/common-fwk/config"
)

func main() {
	cfg, err := config.ValidateConfig(config.NewConfig(
		config.NewServerConfig("0.0.0.0", 8080),
		config.NewSecurityConfig(
			config.NewAuthConfig(
				config.NewJWTConfig("secret", "my-app", 15),
				config.NewCookieConfig("session", "example.com", true, true, "Lax"),
				config.NewLoginConfig("admin@example.com"),
				config.NewOAuth2Config(nil),
			),
		),
		// logging omitted — variadic, optional
	))
	if err != nil {
		log.Fatal(err)
	}

	application := app.NewApplication().
		UseConfig(cfg).
		UseServer()

	application, err = application.UseServerSecurityFromConfig()
	if err != nil {
		log.Fatal(err)
	}

	_ = application.RegisterGET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello"})
	})

	log.Fatal(application.Run())
}
```

### Example 2 — config/viper Load

```go
import viperloader "github.com/matiasmartin-labs/common-fwk/config/viper"

// Load with defaults (reads config.yaml in working directory)
cfg, err := viperloader.Load(viperloader.DefaultOptions())
if err != nil {
	return err
}

// Custom path and env prefix
opts := viperloader.DefaultOptions()
opts.ConfigPath = "/etc/myapp/config.yaml"
opts.EnvPrefix = "APP"
cfg, err = viperloader.Load(opts)
```

### Example 3 — RS256 config constructors

```go
import "github.com/matiasmartin-labs/common-fwk/config"

// Auto-generate RSA key pair (dev/test)
rs256 := config.NewRS256GeneratedConfig("my-key-id")

// Load public PEM (verify-only, production ingress)
rs256pub := config.NewRS256PublicPEMConfig("my-key-id", publicPEM)

// Load private PEM (sign + verify)
rs256priv := config.NewRS256PrivatePEMConfig("my-key-id", privatePEM)

_, _, _ = rs256, rs256pub, rs256priv
```

### Example 4 — security/keys RSA resolvers

```go
import (
	securityjwt "github.com/matiasmartin-labs/common-fwk/security/jwt"
	"github.com/matiasmartin-labs/common-fwk/security/keys"
)

// RSA private key resolver (sign + verify)
resolver := keys.NewRSAResolver(privateKey, "my-key-id")

// RSA public key resolver (verify-only)
resolverPub := keys.NewRSAPublicKeyResolver(publicKey, "my-key-id")

// HMAC static resolver
resolverHMAC := keys.NewStaticResolver(
	&keys.Key{Method: "HS256", Verify: []byte("secret")},
	nil,
)

validator, err := securityjwt.NewValidator(securityjwt.Options{
	Methods:  []string{"RS256"},
	Issuer:   "my-app",
	Resolver: resolver,
})
if err != nil {
	return err
}
_ = resolverPub
_ = resolverHMAC
_ = validator
```

### Example 5 — http/gin middleware options

```go
import (
	"github.com/gin-gonic/gin"
	httpgin "github.com/matiasmartin-labs/common-fwk/http/gin"
)

r := gin.New()

// Full options
r.Use(httpgin.NewAuthMiddleware(validator,
	httpgin.WithAuthEnabled(true),
	httpgin.WithHeaderName("Authorization"),
	httpgin.WithCookieName("session"),
	httpgin.WithContextKey("claims"),
))

// Disable auth conditionally (e.g. in dev mode)
r.Use(httpgin.NewAuthMiddleware(validator, httpgin.WithAuthEnabled(false)))
```

### Example 6 — logging and RSA accessors

```go
// Get a named logger from app
logger, err := application.GetLogger("my-component")
if err != nil {
	log.Fatal(err)
}
_ = logger // implements logging.Logger

// Check security readiness
if application.IsSecurityReady() {
	// Access RSA key material loaded by UseServerSecurityFromConfig
	privKey := application.GetRSAPrivateKey() // *rsa.PrivateKey
	pubKey := application.GetRSAPublicKey()   // *rsa.PublicKey
	keyID := application.GetRSAKeyID()        // string
	_, _, _ = privKey, pubKey, keyID
}
```

## Commands

```bash
go test ./...
```

```bash
go test ./... -run TestAuthMiddleware
```

## Resources

- `README.md`
- `app/doc.go`
- `config/doc.go`
- `config/viper/doc.go`
- `errors/doc.go`
- `http/gin/doc.go`
- `security/doc.go`
- `security/claims/doc.go`
- `security/jwt/doc.go`
- `security/keys/doc.go`
