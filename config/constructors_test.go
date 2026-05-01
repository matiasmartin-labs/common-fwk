package config

import (
	"reflect"
	"testing"
	"time"
)

func TestNewServerConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		host     string
		port     int
		limits   ServerRuntimeLimits
		expected ServerConfig
	}{
		{
			name: "uses defaults when zero values are provided",
			expected: ServerConfig{
				Host:           defaultServerHost,
				Port:           defaultServerPort,
				ReadTimeout:    defaultServerReadTimeout,
				WriteTimeout:   defaultServerWriteTimeout,
				MaxHeaderBytes: defaultServerMaxHeaderBytes,
			},
		},
		{
			name: "keeps explicit values",
			host: "0.0.0.0",
			port: 9090,
			limits: ServerRuntimeLimits{
				ReadTimeout:    3 * time.Second,
				WriteTimeout:   7 * time.Second,
				MaxHeaderBytes: 2048,
			},
			expected: ServerConfig{
				Host:           "0.0.0.0",
				Port:           9090,
				ReadTimeout:    3 * time.Second,
				WriteTimeout:   7 * time.Second,
				MaxHeaderBytes: 2048,
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := NewServerConfig(tc.host, tc.port, tc.limits)
			if actual != tc.expected {
				t.Fatalf("unexpected server config: want %+v, got %+v", tc.expected, actual)
			}
		})
	}
}

func TestNewJWTConfig(t *testing.T) {
	t.Parallel()

	actual := NewJWTConfig("secret", "issuer", 0)
	if actual.TTLMinutes != defaultJWTTTLMinutes {
		t.Fatalf("expected default ttl %d, got %d", defaultJWTTTLMinutes, actual.TTLMinutes)
	}

	if actual.Algorithm != JWTAlgorithmHS256 {
		t.Fatalf("expected default algorithm %q, got %q", JWTAlgorithmHS256, actual.Algorithm)
	}

	if actual.RS256 != (RS256Config{}) {
		t.Fatalf("expected RS256 config to default empty, got %+v", actual.RS256)
	}
}

func TestNewRS256ConfigHelpers(t *testing.T) {
	t.Parallel()

	t.Run("generated key source", func(t *testing.T) {
		cfg := NewRS256GeneratedConfig("kid-1")
		if cfg.KeySource != RS256KeySourceGenerated {
			t.Fatalf("expected key source %q, got %q", RS256KeySourceGenerated, cfg.KeySource)
		}
		if cfg.KeyID != "kid-1" {
			t.Fatalf("expected key id kid-1, got %q", cfg.KeyID)
		}
	})

	t.Run("public pem key source", func(t *testing.T) {
		cfg := NewRS256PublicPEMConfig("kid-2", "PUBLIC")
		if cfg.KeySource != RS256KeySourcePublicPEM {
			t.Fatalf("expected key source %q, got %q", RS256KeySourcePublicPEM, cfg.KeySource)
		}
		if cfg.PublicKeyPEM != "PUBLIC" {
			t.Fatalf("expected public pem to be preserved")
		}
	})

	t.Run("private pem key source", func(t *testing.T) {
		cfg := NewRS256PrivatePEMConfig("kid-3", "PRIVATE")
		if cfg.KeySource != RS256KeySourcePrivatePEM {
			t.Fatalf("expected key source %q, got %q", RS256KeySourcePrivatePEM, cfg.KeySource)
		}
		if cfg.PrivateKeyPEM != "PRIVATE" {
			t.Fatalf("expected private pem to be preserved")
		}
	})
}

func TestNewCookieConfig(t *testing.T) {
	t.Parallel()

	actual := NewCookieConfig("", "", true, false, "")

	if actual.Name != defaultCookieName {
		t.Fatalf("expected default cookie name %q, got %q", defaultCookieName, actual.Name)
	}

	if actual.SameSite != defaultCookieSameSite {
		t.Fatalf("expected default same-site %q, got %q", defaultCookieSameSite, actual.SameSite)
	}

	if actual.HTTPOnly {
		t.Fatalf("expected HTTPOnly to preserve explicit false input")
	}
}

func TestNewConfigIsDeterministicAndCopiesDependencies(t *testing.T) {
	t.Parallel()

	providers := map[string]OAuth2ProviderConfig{
		"google": NewOAuth2ProviderConfig(
			"id",
			"secret",
			"https://accounts.example.com/auth",
			"https://accounts.example.com/token",
			"https://app.example.com/callback",
			[]string{"openid", "email"},
		),
	}

	oauth2 := NewOAuth2Config(providers)
	auth := NewAuthConfig(
		NewJWTConfig("jwt-secret", "common-fwk", 15),
		NewCookieConfig("session", "example.com", true, true, "Lax"),
		NewLoginConfig("user@example.com"),
		oauth2,
	)
	security := NewSecurityConfig(auth)
	server := NewServerConfig("127.0.0.1", 8081)

	first := NewConfig(server, security)
	second := NewConfig(server, security)

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("expected deterministic constructor output")
	}

	providers["google"] = OAuth2ProviderConfig{}
	if first.Security.Auth.OAuth2.Providers["google"].ClientID != "id" {
		t.Fatalf("expected provider map to be copied defensively")
	}

	if !first.Logging.Enabled {
		t.Fatalf("expected logging enabled default to be true")
	}
	if first.Logging.Level != defaultLoggingLevel {
		t.Fatalf("expected default logging level %q, got %q", defaultLoggingLevel, first.Logging.Level)
	}
	if first.Logging.Format != defaultLoggingFormat {
		t.Fatalf("expected default logging format %q, got %q", defaultLoggingFormat, first.Logging.Format)
	}
}

func TestNewLoggingConfigDefaultsAndDefensiveCopy(t *testing.T) {
	t.Parallel()

	enabled := true
	input := map[string]LoggerOverrideConfig{
		"auth": {Enabled: &enabled, Level: "warn"},
	}

	loggingCfg := NewLoggingConfig(true, "", "", input)

	if loggingCfg.Level != defaultLoggingLevel {
		t.Fatalf("expected default level %q, got %q", defaultLoggingLevel, loggingCfg.Level)
	}
	if loggingCfg.Format != defaultLoggingFormat {
		t.Fatalf("expected default format %q, got %q", defaultLoggingFormat, loggingCfg.Format)
	}

	input["auth"] = LoggerOverrideConfig{Level: "error"}
	if loggingCfg.Loggers["auth"].Level != "warn" {
		t.Fatalf("expected logger overrides map to be copied defensively")
	}
}

func TestNewOAuth2ProviderConfigCopiesScopes(t *testing.T) {
	t.Parallel()

	inputScopes := []string{"openid", "email"}
	provider := NewOAuth2ProviderConfig("id", "secret", "auth", "token", "redirect", inputScopes)
	inputScopes[0] = "mutated"

	if provider.Scopes[0] != "openid" {
		t.Fatalf("expected provider scopes to be copied defensively, got %v", provider.Scopes)
	}
}
