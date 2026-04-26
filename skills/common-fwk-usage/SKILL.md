---
name: common-fwk-usage
description: >
  Create or update integration documentation and usage examples for common-fwk.
  Trigger: When writing README quickstarts, architecture overviews, package boundaries,
  or integration examples for config/security/http adapters.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

- Updating README usage docs for `common-fwk`
- Writing quickstart examples that users should follow end-to-end
- Documenting architecture layers and dependency boundaries
- Documenting package responsibilities and explicit non-goals

## Critical Patterns

- Keep examples buildable: include `package main`, imports, and `main()` when snippet scope is full-flow.
- Prefer explicit core API first (`config`, `security/jwt`) and then optional adapters (`config/viper`, `http/gin`).
- Keep boundary clarity explicit: adapters depend on core contracts, never the other way around.
- Preserve non-goals in docs: no app-global singletons, no framework lock-in in core, no remote provider coupling in `security/*`.
- State provider ownership clearly: Google OAuth provider logic stays in consuming apps, outside this framework.

## Documentation Structure

1. One-line project purpose
2. Install step (`go get ...`)
3. Quickstart sequence:
   - explicit core config usage
   - optional `config/viper` facade usage
   - `security/jwt` validator usage
   - `http/gin` middleware integration
4. Layered architecture overview
5. Package responsibilities and boundaries
6. Non-goals

## Code Examples

```go
cfg := config.NewConfig(
	config.NewServerConfig("127.0.0.1", 8080),
	config.NewSecurityConfig(
		config.NewAuthConfig(
			config.NewJWTConfig("secret", "common-fwk", 15),
			config.NewCookieConfig("session", "example.com", true, true, "Lax"),
			config.NewLoginConfig("admin@example.com"),
			config.NewOAuth2Config(nil),
		),
	),
)

validated, err := config.ValidateConfig(cfg)
if err != nil {
	return
}

_ = validated
```

```go
validator, err := securityjwt.NewValidator(securityjwt.Options{
	Methods: []string{"HS256"},
	Issuer:  "common-fwk",
	Resolver: keys.NewStaticResolver(
		&keys.Key{Method: "HS256", Verify: []byte("secret")},
		nil,
	),
})
if err != nil {
	return
}

r := gin.New()
r.Use(ginfwk.NewAuthMiddleware(validator))
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
- `config/doc.go`
- `config/viper/doc.go`
- `security/jwt/doc.go`
- `http/gin/doc.go`
