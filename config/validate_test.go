package config

import (
	"errors"
	"testing"
)

func TestValidateConfigValid(t *testing.T) {
	t.Parallel()

	input := validConfigFixture()
	validated, err := ValidateConfig(input)
	if err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}

	if validated.Security.Auth.Login.Email != "owner@example.com" {
		t.Fatalf("expected normalized email to be preserved in lowercase, got %q", validated.Security.Auth.Login.Email)
	}

	if validated.Security.Auth.JWT.Algorithm != JWTAlgorithmHS256 {
		t.Fatalf("expected default JWT algorithm %q, got %q", JWTAlgorithmHS256, validated.Security.Auth.JWT.Algorithm)
	}
}

func TestValidateConfigInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mutate       func(cfg Config) Config
		wantSentinel error
		wantPath     string
	}{
		{
			name: "missing server host",
			mutate: func(cfg Config) Config {
				cfg.Server.Host = ""
				return cfg
			},
			wantSentinel: ErrRequired,
			wantPath:     "server.host",
		},
		{
			name: "server port out of range",
			mutate: func(cfg Config) Config {
				cfg.Server.Port = 0
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "server.port",
		},
		{
			name: "server read timeout zero",
			mutate: func(cfg Config) Config {
				cfg.Server.ReadTimeout = 0
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "server.readTimeout",
		},
		{
			name: "server read timeout negative",
			mutate: func(cfg Config) Config {
				cfg.Server.ReadTimeout = -1
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "server.readTimeout",
		},
		{
			name: "server write timeout zero",
			mutate: func(cfg Config) Config {
				cfg.Server.WriteTimeout = 0
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "server.writeTimeout",
		},
		{
			name: "server write timeout negative",
			mutate: func(cfg Config) Config {
				cfg.Server.WriteTimeout = -1
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "server.writeTimeout",
		},
		{
			name: "server max header bytes zero",
			mutate: func(cfg Config) Config {
				cfg.Server.MaxHeaderBytes = 0
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "server.maxHeaderBytes",
		},
		{
			name: "server max header bytes negative",
			mutate: func(cfg Config) Config {
				cfg.Server.MaxHeaderBytes = -1
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "server.maxHeaderBytes",
		},
		{
			name: "missing jwt secret",
			mutate: func(cfg Config) Config {
				cfg.Security.Auth.JWT.Secret = ""
				return cfg
			},
			wantSentinel: ErrRequired,
			wantPath:     "security.auth.jwt.secret",
		},
		{
			name: "unsupported jwt algorithm",
			mutate: func(cfg Config) Config {
				cfg.Security.Auth.JWT.Algorithm = "HS512"
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "security.auth.jwt.algorithm",
		},
		{
			name: "rs256 missing key id",
			mutate: func(cfg Config) Config {
				cfg.Security.Auth.JWT.Algorithm = JWTAlgorithmRS256
				cfg.Security.Auth.JWT.Secret = ""
				cfg.Security.Auth.JWT.RS256 = NewRS256GeneratedConfig("")
				return cfg
			},
			wantSentinel: ErrRequired,
			wantPath:     "security.auth.jwt.rs256.keyID",
		},
		{
			name: "rs256 missing public pem",
			mutate: func(cfg Config) Config {
				cfg.Security.Auth.JWT.Algorithm = JWTAlgorithmRS256
				cfg.Security.Auth.JWT.Secret = ""
				cfg.Security.Auth.JWT.RS256 = NewRS256PublicPEMConfig("rsa-1", "")
				return cfg
			},
			wantSentinel: ErrRequired,
			wantPath:     "security.auth.jwt.rs256.publicKeyPEM",
		},
		{
			name: "rs256 missing private pem",
			mutate: func(cfg Config) Config {
				cfg.Security.Auth.JWT.Algorithm = JWTAlgorithmRS256
				cfg.Security.Auth.JWT.Secret = ""
				cfg.Security.Auth.JWT.RS256 = NewRS256PrivatePEMConfig("rsa-1", "")
				return cfg
			},
			wantSentinel: ErrRequired,
			wantPath:     "security.auth.jwt.rs256.privateKeyPEM",
		},
		{
			name: "invalid login email",
			mutate: func(cfg Config) Config {
				cfg.Security.Auth.Login.Email = "not-an-email"
				return cfg
			},
			wantSentinel: ErrInvalidEmail,
			wantPath:     "security.auth.login.email",
		},
		{
			name: "missing oauth2 provider client id",
			mutate: func(cfg Config) Config {
				provider := cfg.Security.Auth.OAuth2.Providers["github"]
				provider.ClientID = ""
				cfg.Security.Auth.OAuth2.Providers["github"] = provider
				return cfg
			},
			wantSentinel: ErrRequired,
			wantPath:     "security.auth.oauth2.providers.github.clientID",
		},
		{
			name: "invalid logging format",
			mutate: func(cfg Config) Config {
				cfg.Logging.Format = "pretty"
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "logging.format",
		},
		{
			name: "invalid root logging level",
			mutate: func(cfg Config) Config {
				cfg.Logging.Level = "verbose"
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "logging.level",
		},
		{
			name: "invalid logger key blank",
			mutate: func(cfg Config) Config {
				cfg.Logging.Loggers["   "] = LoggerOverrideConfig{Level: "debug"}
				return cfg
			},
			wantSentinel: ErrRequired,
			wantPath:     "logging.loggers",
		},
		{
			name: "invalid logger key contains whitespace",
			mutate: func(cfg Config) Config {
				cfg.Logging.Loggers["billing api"] = LoggerOverrideConfig{Level: "debug"}
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "logging.loggers",
		},
		{
			name: "invalid per-logger level",
			mutate: func(cfg Config) Config {
				cfg.Logging.Loggers["auth"] = LoggerOverrideConfig{Level: "notice"}
				return cfg
			},
			wantSentinel: ErrOutOfRange,
			wantPath:     "logging.loggers.auth.level",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := tc.mutate(validConfigFixture())
			_, err := ValidateConfig(cfg)
			if err == nil {
				t.Fatalf("expected validation error")
			}

			if !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("expected ErrInvalidConfig wrapper, got %v", err)
			}

			if !errors.Is(err, tc.wantSentinel) {
				t.Fatalf("expected sentinel %v, got %v", tc.wantSentinel, err)
			}

			var vErr *ValidationError
			if !errors.As(err, &vErr) {
				t.Fatalf("expected ValidationError, got %T", err)
			}

			if vErr.Path != tc.wantPath {
				t.Fatalf("expected path %q, got %q", tc.wantPath, vErr.Path)
			}
		})
	}
}

func TestValidateConfigNormalizesLoggingValues(t *testing.T) {
	t.Parallel()

	enabled := true
	input := validConfigFixture()
	input.Logging.Level = " WARN "
	input.Logging.Format = " JSON "
	input.Logging.Loggers = map[string]LoggerOverrideConfig{
		"auth": {Enabled: &enabled, Level: " DEBUG "},
	}

	validated, err := ValidateConfig(input)
	if err != nil {
		t.Fatalf("expected successful validation, got: %v", err)
	}

	if validated.Logging.Level != "warn" {
		t.Fatalf("expected normalized logging level warn, got %q", validated.Logging.Level)
	}
	if validated.Logging.Format != "json" {
		t.Fatalf("expected normalized logging format json, got %q", validated.Logging.Format)
	}
	if validated.Logging.Loggers["auth"].Level != "debug" {
		t.Fatalf("expected normalized per-logger level debug, got %q", validated.Logging.Loggers["auth"].Level)
	}
}

func TestValidateConfigDefaultsJWTAlgorithmToHS256(t *testing.T) {
	t.Parallel()

	input := validConfigFixture()
	input.Security.Auth.JWT.Algorithm = ""

	validated, err := ValidateConfig(input)
	if err != nil {
		t.Fatalf("expected successful validation, got: %v", err)
	}

	if validated.Security.Auth.JWT.Algorithm != JWTAlgorithmHS256 {
		t.Fatalf("expected normalized algorithm %q, got %q", JWTAlgorithmHS256, validated.Security.Auth.JWT.Algorithm)
	}
}

func TestValidateConfigNormalizesLoginEmail(t *testing.T) {
	t.Parallel()

	input := validConfigFixture()
	input.Security.Auth.Login.Email = "  ADMIN@Example.COM  "

	validated, err := ValidateConfig(input)
	if err != nil {
		t.Fatalf("expected successful validation, got: %v", err)
	}

	if validated.Security.Auth.Login.Email != "admin@example.com" {
		t.Fatalf("expected normalized email, got %q", validated.Security.Auth.Login.Email)
	}
}

func validConfigFixture() Config {
	return NewConfig(
		NewServerConfig("127.0.0.1", 8080),
		NewSecurityConfig(
			NewAuthConfig(
				NewJWTConfig("secret", "common-fwk", 15),
				NewCookieConfig("session", "example.com", true, true, "Lax"),
				NewLoginConfig("owner@example.com"),
				NewOAuth2Config(map[string]OAuth2ProviderConfig{
					"github": NewOAuth2ProviderConfig(
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
}
