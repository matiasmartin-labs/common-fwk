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
			name: "missing jwt secret",
			mutate: func(cfg Config) Config {
				cfg.Security.Auth.JWT.Secret = ""
				return cfg
			},
			wantSentinel: ErrRequired,
			wantPath:     "security.auth.jwt.secret",
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
