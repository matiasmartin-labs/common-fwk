package viper

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestMappingReturnsTypedErrorForInvalidProviderKey(t *testing.T) {
	t.Parallel()

	_, err := mapRawToCore(rawConfig{
		Server: rawServerConfig{Host: "127.0.0.1", Port: 8080, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, MaxHeaderBytes: 1024},
		Security: rawSecurityConfig{Auth: rawAuthConfig{
			JWT:    rawJWTConfig{Secret: "secret", Issuer: "issuer", TTLMinutes: 15},
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
			JWT:    rawJWTConfig{Secret: "secret", Issuer: "issuer", TTLMinutes: 15},
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

	raw.Security.Auth.OAuth2.Providers["github"] = rawOAuth2ProviderConfig{}
	if first.Security.Auth.OAuth2.Providers["github"].ClientID != "id" {
		t.Fatalf("expected mapped providers to be detached from raw source")
	}
}
