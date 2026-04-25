# common-fwk
Common framework

## Config core usage

```go
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
		// handle classified validation failures
	}
}

// validated.Security.Auth.Login.Email == "admin@example.com"
```
