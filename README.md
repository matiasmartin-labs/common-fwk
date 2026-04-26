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

## Optional Viper adapter usage

```go
loaded, err := viper.Load(viper.Options{
	ConfigPath:  "./config/app.yaml",
	ConfigType:  "",          // infer from extension when empty
	EnvPrefix:   "COMMON_FWK", // default when omitted
	EnvOverride: true,          // env values override file values
	ExpandEnv:   true,          // expand ${VAR} placeholders from env snapshot
})
if err != nil {
	var loadErr *viper.LoadError
	var decodeErr *viper.DecodeError
	var mapErr *viper.MappingError
	var validateErr *viper.ValidationError

	switch {
	case errors.As(err, &loadErr):
		// file access or option/application failure
	case errors.As(err, &decodeErr):
		// parse/unmarshal failure
	case errors.As(err, &mapErr):
		// explicit raw -> core mapping failure
	case errors.As(err, &validateErr):
		// wraps config.ValidateConfig failures
	}

	if errors.Is(err, config.ErrInvalidConfig) {
		// preserved core classification through ValidationError wrapping
	}
}

_ = loaded
```

Behavior notes:
- `EnvOverride=false` keeps file values as source of truth.
- `EnvOverride=true` applies env values using `EnvPrefix` and deterministic keys (for example `COMMON_FWK_SECURITY_AUTH_JWT_SECRET`).
- `ExpandEnv=false` preserves `${VAR}` placeholders.
- `ExpandEnv=true` expands placeholders against a per-load env snapshot deterministically.

## Security core JWT validation

`security/*` provides a framework-agnostic JWT validation core with deterministic dependencies.

```go
jwtCfg := config.NewJWTConfig("secret", "common-fwk", 15)
compat := securityjwt.FromConfigJWT(jwtCfg)

validator, err := securityjwt.NewValidator(compat.Options)
if err != nil {
	// invalid validator options (for example missing resolver)
}

validatedClaims, err := validator.Validate(ctx, rawToken)
if err != nil {
	if errors.Is(err, securityjwt.ErrInvalidMethod) {
		// token alg is not in allowlist
	}
	if errors.Is(err, securityjwt.ErrExpiredToken) {
		// token exp is before validator clock
	}

	var vErr *securityjwt.ValidationError
	if errors.As(err, &vErr) {
		// inspect stage metadata while preserving sentinel assertability
	}
}

_ = validatedClaims
_ = compat.TokenTTL // token issuing concern, not validator runtime enforcement
```

Package boundaries:
- `security/claims`: standard claim model and audience normalization helpers.
- `security/keys`: deterministic key resolution contracts (`kid` + default key fallback).
- `security/jwt`: validator options, typed/sentinel error taxonomy, validation flow, config compatibility helper.

Explicit non-goals of this security core:
- no Gin middleware coupling
- no app-global singleton coupling
- no OAuth/JWKS provider adapter/network fetch coupling
