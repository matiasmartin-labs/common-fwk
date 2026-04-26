package viper

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/matiasmartin-labs/common-fwk/config"
)

func TestLoadSuccessAndDeterminism(t *testing.T) {
	t.Parallel()

	configPath := writeTestConfig(t, "valid.yaml", `
server:
  host: "127.0.0.1"
  port: 8080
security:
  auth:
    jwt:
      secret: "secret"
      issuer: "common-fwk"
      ttlMinutes: 15
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      httpOnly: true
      sameSite: "Lax"
    login:
      email: "OWNER@Example.com"
    oauth2:
      providers:
        github:
          clientID: "id"
          clientSecret: "secret"
          authURL: "https://github.com/login/oauth/authorize"
          tokenURL: "https://github.com/login/oauth/access_token"
          redirectURL: "https://app.example.com/auth/github/callback"
          scopes: ["read:user", "user:email"]
`)

	first, err := Load(Options{ConfigPath: configPath})
	if err != nil {
		t.Fatalf("expected successful load, got error: %v", err)
	}

	second, err := Load(Options{ConfigPath: configPath})
	if err != nil {
		t.Fatalf("expected repeated load to succeed, got error: %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("expected deterministic output for identical inputs")
	}

	if first.Security.Auth.Login.Email != "owner@example.com" {
		t.Fatalf("expected normalized email from core validation, got %q", first.Security.Auth.Login.Email)
	}
}

func TestLoadFailureTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		opts        Options
		setup       func(t *testing.T) Options
		assertError func(t *testing.T, err error)
	}{
		{
			name: "missing file returns load error",
			setup: func(t *testing.T) Options {
				t.Helper()
				return Options{ConfigPath: filepath.Join(t.TempDir(), "missing.yaml")}
			},
			assertError: func(t *testing.T, err error) {
				t.Helper()
				var loadErr *LoadError
				if !errors.As(err, &loadErr) {
					t.Fatalf("expected LoadError, got %T", err)
				}
			},
		},
		{
			name: "malformed content returns decode error",
			setup: func(t *testing.T) Options {
				t.Helper()
				path := writeTestConfig(t, "broken.yaml", "server: [unclosed")
				return Options{ConfigPath: path}
			},
			assertError: func(t *testing.T, err error) {
				t.Helper()
				var decodeErr *DecodeError
				if !errors.As(err, &decodeErr) {
					t.Fatalf("expected DecodeError, got %T", err)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := tc.opts
			if tc.setup != nil {
				opts = tc.setup(t)
			}

			_, err := Load(opts)
			if err == nil {
				t.Fatalf("expected load error")
			}

			tc.assertError(t, err)
		})
	}
}

func TestLoadEnvOverrideSemantics(t *testing.T) {
	path := writeTestConfig(t, "env-override.yaml", `
server:
  host: "127.0.0.1"
  port: 8080
security:
  auth:
    jwt:
      secret: "file-secret"
      issuer: "common-fwk"
      ttlMinutes: 15
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      httpOnly: true
      sameSite: "Lax"
    login:
      email: "owner@example.com"
    oauth2:
      providers: {}
`)

	t.Setenv(defaultEnvPrefix+"_SECURITY_AUTH_JWT_SECRET", "env-secret")

	withoutOverride, err := Load(Options{ConfigPath: path, EnvOverride: false})
	if err != nil {
		t.Fatalf("unexpected load error with EnvOverride=false: %v", err)
	}

	withOverride, err := Load(Options{ConfigPath: path, EnvOverride: true})
	if err != nil {
		t.Fatalf("unexpected load error with EnvOverride=true: %v", err)
	}

	if withoutOverride.Security.Auth.JWT.Secret != "file-secret" {
		t.Fatalf("expected file value when EnvOverride=false, got %q", withoutOverride.Security.Auth.JWT.Secret)
	}

	if withOverride.Security.Auth.JWT.Secret != "env-secret" {
		t.Fatalf("expected env value when EnvOverride=true, got %q", withOverride.Security.Auth.JWT.Secret)
	}
}

func TestLoadExpandEnvDeterminism(t *testing.T) {
	path := writeTestConfig(t, "expand.yaml", `
server:
  host: "${APP_HOST}"
  port: 8080
security:
  auth:
    jwt:
      secret: "secret"
      issuer: "common-fwk"
      ttlMinutes: 15
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      httpOnly: true
      sameSite: "Lax"
    login:
      email: "owner@example.com"
    oauth2:
      providers: {}
`)

	t.Setenv("APP_HOST", "10.10.10.10")

	withoutExpansion, err := Load(Options{ConfigPath: path, ExpandEnv: false})
	if err != nil {
		t.Fatalf("unexpected load error with ExpandEnv=false: %v", err)
	}

	withExpansionFirst, err := Load(Options{ConfigPath: path, ExpandEnv: true})
	if err != nil {
		t.Fatalf("unexpected load error with ExpandEnv=true: %v", err)
	}

	withExpansionSecond, err := Load(Options{ConfigPath: path, ExpandEnv: true})
	if err != nil {
		t.Fatalf("unexpected repeated load error with ExpandEnv=true: %v", err)
	}

	if withoutExpansion.Server.Host != "${APP_HOST}" {
		t.Fatalf("expected placeholder to remain when ExpandEnv=false, got %q", withoutExpansion.Server.Host)
	}

	if withExpansionFirst.Server.Host != "10.10.10.10" {
		t.Fatalf("expected expanded host when ExpandEnv=true, got %q", withExpansionFirst.Server.Host)
	}

	if !reflect.DeepEqual(withExpansionFirst, withExpansionSecond) {
		t.Fatalf("expected deterministic expanded output for fixed env snapshot")
	}
}

func TestLoadWrapsCoreValidation(t *testing.T) {
	t.Parallel()

	path := writeTestConfig(t, "invalid-core.yaml", `
server:
  host: "127.0.0.1"
  port: 8080
security:
  auth:
    jwt:
      secret: ""
      issuer: "common-fwk"
      ttlMinutes: 15
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      httpOnly: true
      sameSite: "Lax"
    login:
      email: "owner@example.com"
    oauth2:
      providers: {}
`)

	_, err := Load(Options{ConfigPath: path})
	if err == nil {
		t.Fatalf("expected validation failure")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if !errors.Is(err, config.ErrInvalidConfig) {
		t.Fatalf("expected wrapped error to preserve config.ErrInvalidConfig")
	}

	if !errors.Is(err, config.ErrRequired) {
		t.Fatalf("expected wrapped error to preserve config.ErrRequired")
	}

	var coreValidation *config.ValidationError
	if !errors.As(err, &coreValidation) {
		t.Fatalf("expected core ValidationError to remain assertable")
	}
}

func writeTestConfig(t *testing.T, name, contents string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write test config: %v", err)
	}

	return path
}
