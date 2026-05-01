---
title: Quickstart
parent: Getting Started
nav_order: 2
---

# Quickstart

This example wires config, JWT HS256 security, and an HTTP server in a single file.

## 1. Config file (`config.yaml`)

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
      algorithm: HS256
      secret: my-secret
      issuer: my-service
      ttl-minutes: 60

logging:
  enabled: true
  level: info
  format: json
```

## 2. Bootstrap (`main.go`)

```go
package main

import (
    "log"

    "github.com/matiasmartin-labs/common-fwk/app"
    viperconfig "github.com/matiasmartin-labs/common-fwk/config/viper"
    "github.com/matiasmartin-labs/common-fwk/security/jwt"
    "github.com/matiasmartin-labs/common-fwk/security/keys"
    httpgin "github.com/matiasmartin-labs/common-fwk/http/gin"
    "github.com/gin-gonic/gin"
)

func main() {
    cfg, err := viperconfig.Load("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    resolver := keys.NewStaticResolver([]byte(cfg.Security.Auth.JWT.Secret))
    validator, err := jwt.NewValidator(cfg.Security.Auth.JWT, resolver)
    if err != nil {
        log.Fatal(err)
    }

    application := app.NewApplication()
    application.UseConfig(cfg)
    application.UseServer()
    application.UseServerSecurity(validator)

    application.RegisterGET("/healthz", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    application.RegisterProtectedGET("/me", httpgin.NewAuthMiddleware(validator), func(c *gin.Context) {
        c.JSON(200, gin.H{"user": "authenticated"})
    })

    if err := application.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## 3. Run

```bash
go run main.go
```

## Next steps

- Read the [Architecture](../architecture/) section for detailed contracts.
- See the [Migration Guide](../migration/auth-provider-ms-v0.1.0/) if migrating from `auth-provider-ms/pkg`.
