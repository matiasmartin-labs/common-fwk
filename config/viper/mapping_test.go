package viper

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/matiasmartin-labs/common-fwk/config"
)

func TestMappingReturnsTypedErrorForInvalidProviderKey(t *testing.T) {
	t.Parallel()

	_, err := mapRawToCore(rawConfig{
		Server: rawServerConfig{Host: "127.0.0.1", Port: 8080, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, MaxHeaderBytes: 1024},
		Security: rawSecurityConfig{Auth: rawAuthConfig{
			JWT:    rawJWTConfig{Algorithm: "HS256", Secret: "secret", Issuer: "issuer", TTLMinutes: 15},
			Cookie: rawCookieConfig{Name: "session", SameSite: "Lax"},
			Login:  rawLoginConfig{Email: "owner@example.com"},
			OAuth2: rawOAuth2Config{Providers: map[string]rawOAuth2ProviderConfig{
				"   ": {
					ClientID:     "id",
					ClientSecret: "secret",
					AuthURL:      "https://accounts.example.com/auth",
					TokenURL:     "https://accounts.example.com/token",
					RedirectURL:  "https://app.example.com/callback",
				},
			}},
		}},
	})
	if err == nil {
		t.Fatalf("expected mapping error")
	}

	var mappingErr *MappingError
	if !errors.As(err, &mappingErr) {
		t.Fatalf("expected MappingError, got %T", err)
	}

	if mappingErr.Path != "security.auth.oauth2.providers" {
		t.Fatalf("expected providers path, got %q", mappingErr.Path)
	}

	if !errors.Is(err, errEmptyProviderKey) {
		t.Fatalf("expected provider-key classification to be preserved")
	}
}

func TestMappingDeterministicAndDefensiveCopies(t *testing.T) {
	t.Parallel()

	raw := rawConfig{
		Server: rawServerConfig{Host: "127.0.0.1", Port: 8080, ReadTimeout: 9 * time.Second, WriteTimeout: 11 * time.Second, MaxHeaderBytes: 2048},
		Security: rawSecurityConfig{Auth: rawAuthConfig{
			JWT: rawJWTConfig{
				Algorithm:         "RS256",
				Issuer:            "issuer",
				TTLMinutes:        15,
				RS256KeySource:    "public-pem",
				RS256KeyID:        "rsa-key",
				RS256PublicKeyPEM: "PUBLIC",
			},
			Cookie: rawCookieConfig{Name: "session", Domain: "example.com", Secure: true, HTTPOnly: true, SameSite: "Lax"},
			Login:  rawLoginConfig{Email: "owner@example.com"},
			OAuth2: rawOAuth2Config{Providers: map[string]rawOAuth2ProviderConfig{
				"github": {
					ClientID:     "id",
					ClientSecret: "secret",
					AuthURL:      "https://github.com/login/oauth/authorize",
					TokenURL:     "https://github.com/login/oauth/access_token",
					RedirectURL:  "https://app.example.com/auth/github/callback",
					Scopes:       []string{"read:user", "user:email"},
				},
			}},
		}},
	}

	first, err := mapRawToCore(raw)
	if err != nil {
		t.Fatalf("unexpected first mapping error: %v", err)
	}

	second, err := mapRawToCore(raw)
	if err != nil {
		t.Fatalf("unexpected second mapping error: %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("expected deterministic mapping outputs")
	}

	if first.Server.ReadTimeout != 9*time.Second {
		t.Fatalf("expected read timeout to be mapped, got %s", first.Server.ReadTimeout)
	}
	if first.Server.WriteTimeout != 11*time.Second {
		t.Fatalf("expected write timeout to be mapped, got %s", first.Server.WriteTimeout)
	}
	if first.Server.MaxHeaderBytes != 2048 {
		t.Fatalf("expected max header bytes to be mapped, got %d", first.Server.MaxHeaderBytes)
	}

	if first.Security.Auth.JWT.Algorithm != config.JWTAlgorithmRS256 {
		t.Fatalf("expected algorithm RS256, got %q", first.Security.Auth.JWT.Algorithm)
	}
	if first.Security.Auth.JWT.RS256.KeySource != config.RS256KeySourcePublicPEM {
		t.Fatalf("expected key source public-pem, got %q", first.Security.Auth.JWT.RS256.KeySource)
	}
	if first.Security.Auth.JWT.RS256.KeyID != "rsa-key" {
		t.Fatalf("expected key id rsa-key, got %q", first.Security.Auth.JWT.RS256.KeyID)
	}

	if !first.Logging.Enabled {
		t.Fatalf("expected logging enabled default true")
	}
	if first.Logging.Level != "info" {
		t.Fatalf("expected default logging level info, got %q", first.Logging.Level)
	}
	if first.Logging.Format != "json" {
		t.Fatalf("expected default logging format json, got %q", first.Logging.Format)
	}

	raw.Security.Auth.OAuth2.Providers["github"] = rawOAuth2ProviderConfig{}
	if first.Security.Auth.OAuth2.Providers["github"].ClientID != "id" {
		t.Fatalf("expected mapped providers to be detached from raw source")
	}
}

func TestMappingIncludesLoggingRootAndPerLoggerOverrides(t *testing.T) {
	t.Parallel()

	enabled := true
	raw := rawConfig{
		Server: rawServerConfig{Host: "127.0.0.1", Port: 8080, ReadTimeout: 5 * time.Second, WriteTimeout: 5 * time.Second, MaxHeaderBytes: 1024},
		Security: rawSecurityConfig{Auth: rawAuthConfig{
			JWT:    rawJWTConfig{Algorithm: "HS256", Secret: "secret", Issuer: "issuer", TTLMinutes: 15},
			Cookie: rawCookieConfig{Name: "session", SameSite: "Lax"},
			Login:  rawLoginConfig{Email: "owner@example.com"},
			OAuth2: rawOAuth2Config{Providers: map[string]rawOAuth2ProviderConfig{}},
		}},
		Logging: rawLoggingConfig{
			Enabled: &enabled,
			Level:   "warn",
			Format:  "text",
			Loggers: map[string]rawLoggerOverrideConfig{
				"auth": {Level: "debug"},
			},
		},
	}

	mapped, err := mapRawToCore(raw)
	if err != nil {
		t.Fatalf("unexpected mapping error: %v", err)
	}

	if !mapped.Logging.Enabled {
		t.Fatalf("expected mapped logging enabled true")
	}
	if mapped.Logging.Level != "warn" {
		t.Fatalf("expected mapped logging level warn, got %q", mapped.Logging.Level)
	}
	if mapped.Logging.Format != "text" {
		t.Fatalf("expected mapped logging format text, got %q", mapped.Logging.Format)
	}
	if mapped.Logging.Loggers["auth"].Level != "debug" {
		t.Fatalf("expected per-logger level override to map")
	}
}
