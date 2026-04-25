package config

import (
	"reflect"
	"testing"
)

func TestNewServerConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		host     string
		port     int
		expected ServerConfig
	}{
		{
			name: "uses defaults when zero values are provided",
			expected: ServerConfig{
				Host: defaultServerHost,
				Port: defaultServerPort,
			},
		},
		{
			name: "keeps explicit values",
			host: "0.0.0.0",
			port: 9090,
			expected: ServerConfig{
				Host: "0.0.0.0",
				Port: 9090,
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := NewServerConfig(tc.host, tc.port)
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
