---
title: Gin Middleware
parent: Architecture
nav_order: 5
---

# Gin Auth Middleware (`http/gin`)

**Import**: `github.com/matiasmartin-labs/common-fwk/http/gin`

## Purpose

Gin HTTP middleware that authenticates JWT Bearer tokens using `security.Validator`,
returns standardized error responses, and injects validated claims into request context.

## Factory

```go
middleware := httpgin.NewAuthMiddleware(validator, opts...)
```

## Options

| Option | Default | Description |
|---|---|---|
| `WithHeaderName(name)` | `Authorization` | Header to extract Bearer token from |
| `WithCookieName(name)` | `token` | Cookie fallback name |
| `WithAuthEnabled(bool)` | `true` | Set `false` to bypass auth (testing/dev) |
| `WithContextKey(key)` | internal default | Key to inject claims into `gin.Context` |

## Token Extraction Precedence

1. Configured header (strip `Bearer ` prefix)
2. Configured cookie (fallback when header is absent)
3. Neither present → `401` with `auth_token_missing`

## Error Response Contract

All auth failures return `HTTP 401` with JSON body:

```json
{ "code": "auth_token_missing", "message": "missing authentication token" }
{ "code": "auth_token_invalid", "message": "invalid or expired token" }
```

| Failure | Code |
|---|---|
| No token found | `auth_token_missing` |
| Malformed / invalid / expired / bad issuer or audience | `auth_token_invalid` |

## Exported Constants

```go
httpgin.MsgTokenMissing  // "missing authentication token"
httpgin.MsgTokenInvalid  // "invalid or expired token"
```

Use these to assert error messages in tests without magic strings.

## Exported Error Codes (from `errors` package)

```go
import fwkerrors "github.com/matiasmartin-labs/common-fwk/errors"

fwkerrors.CodeTokenMissing  // "auth_token_missing"
fwkerrors.CodeTokenInvalid  // "auth_token_invalid"
```

## Claims Injection

On success, validated `claims.Claims` is injected into `gin.Context` under the configured key.
Downstream handlers retrieve it via the standard Gin context API.

## Usage Example

```go
application.RegisterProtectedGET("/profile",
    httpgin.NewAuthMiddleware(validator),
    func(c *gin.Context) {
        c.JSON(200, gin.H{"ok": true})
    },
)
```
