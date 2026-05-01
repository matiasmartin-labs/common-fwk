package viper

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/matiasmartin-labs/common-fwk/config"
)

func TestLoadSuccessAndDeterminism(t *testing.T) {
	t.Parallel()

	configPath := writeTestConfig(t, "valid.yaml", `
server:
  host: "127.0.0.1"
  port: 8080
  read-timeout: 4s
  write-timeout: 6s
  max-header-bytes: 8192
security:
  auth:
    jwt:
      secret: "secret"
      issuer: "common-fwk"
      ttl-minutes: 15
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      http-only: true
      same-site: "Lax"
    login:
      email: "OWNER@Example.com"
    oauth2:
      providers:
        github:
          client-id: "id"
          client-secret: "secret"
          auth-url: "https://github.com/login/oauth/authorize"
          token-url: "https://github.com/login/oauth/access_token"
          redirect-url: "https://app.example.com/auth/github/callback"
          scopes: ["read:user", "user:email"]
logging:
  enabled: true
  level: info
  format: json
  loggers:
    auth:
      enabled: true
      level: debug
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

	if first.Server.ReadTimeout != 4*time.Second {
		t.Fatalf("expected read-timeout mapping, got %s", first.Server.ReadTimeout)
	}

	if first.Server.WriteTimeout != 6*time.Second {
		t.Fatalf("expected write-timeout mapping, got %s", first.Server.WriteTimeout)
	}

	if first.Server.MaxHeaderBytes != 8192 {
		t.Fatalf("expected max-header-bytes mapping, got %d", first.Server.MaxHeaderBytes)
	}

	if first.Security.Auth.JWT.Algorithm != config.JWTAlgorithmHS256 {
		t.Fatalf("expected default algorithm HS256, got %q", first.Security.Auth.JWT.Algorithm)
	}

	if first.Logging.Level != "info" || first.Logging.Format != "json" {
		t.Fatalf("expected logging root mapping to be preserved")
	}
	if first.Logging.Loggers["auth"].Level != "debug" {
		t.Fatalf("expected per-logger logging override mapping")
	}
}

func TestLoadRS256CanonicalAndLegacyMapping(t *testing.T) {
	t.Parallel()

	path := writeTestConfig(t, "rs256-mixed.yaml", `
server:
  host: "127.0.0.1"
  port: 8080
security:
  auth:
    jwt:
      algorithm: "RS256"
      secret: "ignored-for-rs256"
      issuer: "common-fwk"
      ttl-minutes: 20
      rs256-key-source: "public-pem"
      rs256-key-id: "canonical-key"
      rs256-public-key-pem: "CANONICAL_PUBLIC"
      rs256KeyID: "legacy-key"
      rs256PublicKeyPEM: "LEGACY_PUBLIC"
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      http-only: true
      same-site: "Lax"
    login:
      email: "owner@example.com"
    oauth2:
      providers: {}
logging:
  enabled: true
  level: info
  format: json
  loggers:
    auth:
      level: info
`)

	first, err := Load(Options{ConfigPath: path})
	if err != nil {
		t.Fatalf("expected RS256 load success, got error: %v", err)
	}

	second, err := Load(Options{ConfigPath: path})
	if err != nil {
		t.Fatalf("expected repeated RS256 load success, got error: %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("expected deterministic RS256 mapping output")
	}

	jwtCfg := first.Security.Auth.JWT
	if jwtCfg.Algorithm != config.JWTAlgorithmRS256 {
		t.Fatalf("expected RS256 algorithm, got %q", jwtCfg.Algorithm)
	}
	if jwtCfg.RS256.KeySource != config.RS256KeySourcePublicPEM {
		t.Fatalf("expected key source public-pem, got %q", jwtCfg.RS256.KeySource)
	}
	if jwtCfg.RS256.KeyID != "canonical-key" {
		t.Fatalf("expected canonical key id to win, got %q", jwtCfg.RS256.KeyID)
	}
	if jwtCfg.RS256.PublicKeyPEM != "CANONICAL_PUBLIC" {
		t.Fatalf("expected canonical public PEM to win, got %q", jwtCfg.RS256.PublicKeyPEM)
	}
}

func TestLoadRS256EnvOverrides(t *testing.T) {
	path := writeTestConfig(t, "rs256-env.yaml", `
server:
  host: "127.0.0.1"
  port: 8080
security:
  auth:
    jwt:
      algorithm: "HS256"
      secret: "file-secret"
      issuer: "common-fwk"
      ttl-minutes: 15
      rs256-key-source: "generated"
      rs256-key-id: "file-key-id"
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      http-only: true
      same-site: "Lax"
    login:
      email: "owner@example.com"
    oauth2:
      providers: {}
logging:
  enabled: true
  level: info
  format: json
`)

	t.Setenv(defaultEnvPrefix+"_SECURITY_AUTH_JWT_ALGORITHM", "RS256")
	t.Setenv(defaultEnvPrefix+"_SECURITY_AUTH_JWT_RS256_KEY_SOURCE", "private-pem")
	t.Setenv(defaultEnvPrefix+"_SECURITY_AUTH_JWT_RS256_KEY_ID", "env-key-id")
	t.Setenv(defaultEnvPrefix+"_SECURITY_AUTH_JWT_RS256_PRIVATE_KEY_PEM", "ENV_PRIVATE")

	cfg, err := Load(Options{ConfigPath: path, EnvOverride: true})
	if err != nil {
		t.Fatalf("expected RS256 env override load success, got error: %v", err)
	}

	jwtCfg := cfg.Security.Auth.JWT
	if jwtCfg.Algorithm != config.JWTAlgorithmRS256 {
		t.Fatalf("expected env algorithm RS256, got %q", jwtCfg.Algorithm)
	}
	if jwtCfg.RS256.KeySource != config.RS256KeySourcePrivatePEM {
		t.Fatalf("expected env key source private-pem, got %q", jwtCfg.RS256.KeySource)
	}
	if jwtCfg.RS256.KeyID != "env-key-id" {
		t.Fatalf("expected env key id, got %q", jwtCfg.RS256.KeyID)
	}
	if jwtCfg.RS256.PrivateKeyPEM != "ENV_PRIVATE" {
		t.Fatalf("expected env private key pem, got %q", jwtCfg.RS256.PrivateKeyPEM)
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
  read-timeout: 2s
  write-timeout: 3s
  max-header-bytes: 4096
security:
  auth:
    jwt:
      secret: "file-secret"
      issuer: "common-fwk"
      ttl-minutes: 15
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      http-only: true
      same-site: "Lax"
    login:
      email: "owner@example.com"
    oauth2:
      providers: {}
logging:
  enabled: true
  level: info
  format: json
`)

	t.Setenv(defaultEnvPrefix+"_SECURITY_AUTH_JWT_SECRET", "env-secret")
	t.Setenv(defaultEnvPrefix+"_SERVER_READ_TIMEOUT", "8s")
	t.Setenv(defaultEnvPrefix+"_SERVER_WRITE_TIMEOUT", "9s")
	t.Setenv(defaultEnvPrefix+"_SERVER_MAX_HEADER_BYTES", "16384")
	t.Setenv(defaultEnvPrefix+"_LOGGING_LEVEL", "warn")
	t.Setenv(defaultEnvPrefix+"_LOGGING_LOGGERS_AUTH_LEVEL", "error")

	withoutOverride, err := Load(Options{ConfigPath: path, EnvOverride: false})
	if err != nil {
		t.Fatalf("unexpected load error with EnvOverride=false: %v", err)
	}

	withOverride, err := Load(Options{ConfigPath: path, EnvOverride: true})
	if err != nil {
		t.Fatalf("unexpected load error with EnvOverride=true: %v", err)
	}

	withOverrideSecond, err := Load(Options{ConfigPath: path, EnvOverride: true})
	if err != nil {
		t.Fatalf("unexpected repeated load error with EnvOverride=true: %v", err)
	}

	if withoutOverride.Security.Auth.JWT.Secret != "file-secret" {
		t.Fatalf("expected file value when EnvOverride=false, got %q", withoutOverride.Security.Auth.JWT.Secret)
	}

	if withOverride.Security.Auth.JWT.Secret != "env-secret" {
		t.Fatalf("expected env value when EnvOverride=true, got %q", withOverride.Security.Auth.JWT.Secret)
	}

	if withoutOverride.Server.ReadTimeout != 2*time.Second || withoutOverride.Server.WriteTimeout != 3*time.Second || withoutOverride.Server.MaxHeaderBytes != 4096 {
		t.Fatalf("expected file server runtime limits when EnvOverride=false")
	}

	if withOverride.Server.ReadTimeout != 8*time.Second || withOverride.Server.WriteTimeout != 9*time.Second || withOverride.Server.MaxHeaderBytes != 16384 {
		t.Fatalf("expected env server runtime limits when EnvOverride=true")
	}

	if withoutOverride.Logging.Level != "info" {
		t.Fatalf("expected file logging level when EnvOverride=false, got %q", withoutOverride.Logging.Level)
	}
	if withOverride.Logging.Level != "warn" {
		t.Fatalf("expected env logging.level when EnvOverride=true, got %q", withOverride.Logging.Level)
	}
	if withOverride.Logging.Loggers["auth"].Level != "error" {
		t.Fatalf("expected env logger override level when EnvOverride=true")
	}

	if !reflect.DeepEqual(withOverride, withOverrideSecond) {
		t.Fatalf("expected deterministic output for identical env snapshot when EnvOverride=true")
	}
}

func TestLoadEnvOverrideTypedFailuresForServerRuntimeLimits(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		wantSubstr string
	}{
		{
			name:       "invalid server read timeout format",
			envKey:     defaultEnvPrefix + "_SERVER_READ_TIMEOUT",
			envValue:   "not-a-duration",
			wantSubstr: "SERVER_READ_TIMEOUT",
		},
		{
			name:       "invalid server write timeout format",
			envKey:     defaultEnvPrefix + "_SERVER_WRITE_TIMEOUT",
			envValue:   "still-not-a-duration",
			wantSubstr: "SERVER_WRITE_TIMEOUT",
		},
		{
			name:       "invalid max header bytes format",
			envKey:     defaultEnvPrefix + "_SERVER_MAX_HEADER_BYTES",
			envValue:   "NaN",
			wantSubstr: "SERVER_MAX_HEADER_BYTES",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			path := writeTestConfig(t, "env-typed-failure.yaml", `
server:
  host: "127.0.0.1"
  port: 8080
  read-timeout: 2s
  write-timeout: 3s
  max-header-bytes: 4096
security:
  auth:
    jwt:
      secret: "file-secret"
      issuer: "common-fwk"
      ttl-minutes: 15
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      http-only: true
      same-site: "Lax"
    login:
      email: "owner@example.com"
    oauth2:
      providers: {}
logging:
  enabled: true
  level: info
  format: json
`)

			t.Setenv(tc.envKey, tc.envValue)

			_, err := Load(Options{ConfigPath: path, EnvOverride: true})
			if err == nil {
				t.Fatalf("expected load error")
			}

			var loadErr *LoadError
			if !errors.As(err, &loadErr) {
				t.Fatalf("expected LoadError, got %T", err)
			}

			if !strings.Contains(err.Error(), tc.wantSubstr) {
				t.Fatalf("expected error to mention %q, got %v", tc.wantSubstr, err)
			}
		})
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
      ttl-minutes: 15
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      http-only: true
      same-site: "Lax"
    login:
      email: "owner@example.com"
    oauth2:
      providers: {}
logging:
  enabled: true
  level: info
  format: json
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
      ttl-minutes: 15
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      http-only: true
      same-site: "Lax"
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

func TestLoadWrapsCoreValidationForInvalidLoggingValues(t *testing.T) {
	t.Parallel()

	path := writeTestConfig(t, "invalid-logging.yaml", `
server:
  host: "127.0.0.1"
  port: 8080
security:
  auth:
    jwt:
      secret: "secret"
      issuer: "common-fwk"
      ttl-minutes: 15
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      http-only: true
      same-site: "Lax"
    login:
      email: "owner@example.com"
    oauth2:
      providers: {}
logging:
  enabled: true
  level: info
  format: pretty
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

	if !errors.Is(err, config.ErrOutOfRange) {
		t.Fatalf("expected wrapped error to preserve config.ErrOutOfRange")
	}

	var coreValidation *config.ValidationError
	if !errors.As(err, &coreValidation) {
		t.Fatalf("expected core ValidationError to remain assertable")
	}
}

func TestLoadLegacyCamelCaseCompatibility(t *testing.T) {
	t.Parallel()

	path := writeTestConfig(t, "legacy-camel-case.yaml", `
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
      email: "owner@example.com"
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

	cfg, err := Load(Options{ConfigPath: path})
	if err != nil {
		t.Fatalf("expected legacy camelCase keys to remain compatible, got error: %v", err)
	}

	provider := cfg.Security.Auth.OAuth2.Providers["github"]
	if cfg.Security.Auth.JWT.TTLMinutes != 15 {
		t.Fatalf("expected ttlMinutes compatibility mapping, got %d", cfg.Security.Auth.JWT.TTLMinutes)
	}
	if !cfg.Security.Auth.Cookie.HTTPOnly {
		t.Fatalf("expected httpOnly compatibility mapping")
	}
	if cfg.Security.Auth.Cookie.SameSite != "Lax" {
		t.Fatalf("expected sameSite compatibility mapping, got %q", cfg.Security.Auth.Cookie.SameSite)
	}
	if provider.ClientID != "id" {
		t.Fatalf("expected provider clientID compatibility mapping, got %q", provider.ClientID)
	}
}

func TestLoadCanonicalPrecedenceOverLegacyKeys(t *testing.T) {
	t.Parallel()

	path := writeTestConfig(t, "mixed-style.yaml", `
server:
  host: "127.0.0.1"
  port: 8080
security:
  auth:
    jwt:
      secret: "secret"
      issuer: "common-fwk"
      ttl-minutes: 20
      ttlMinutes: 10
    cookie:
      name: "session"
      domain: "example.com"
      secure: true
      http-only: true
      httpOnly: false
      same-site: "Strict"
      sameSite: "Lax"
    login:
      email: "owner@example.com"
    oauth2:
      providers:
        github:
          client-id: "canonical-id"
          clientID: "legacy-id"
          client-secret: "canonical-secret"
          clientSecret: "legacy-secret"
          auth-url: "https://canonical.example.com/auth"
          authURL: "https://legacy.example.com/auth"
          token-url: "https://canonical.example.com/token"
          tokenURL: "https://legacy.example.com/token"
          redirect-url: "https://canonical.example.com/callback"
          redirectURL: "https://legacy.example.com/callback"
          scopes: ["read:user"]
`)

	cfg, err := Load(Options{ConfigPath: path})
	if err != nil {
		t.Fatalf("expected mixed-style config to load, got error: %v", err)
	}

	provider := cfg.Security.Auth.OAuth2.Providers["github"]
	if cfg.Security.Auth.JWT.TTLMinutes != 20 {
		t.Fatalf("expected canonical ttl-minutes to win, got %d", cfg.Security.Auth.JWT.TTLMinutes)
	}
	if !cfg.Security.Auth.Cookie.HTTPOnly {
		t.Fatalf("expected canonical http-only to win")
	}
	if cfg.Security.Auth.Cookie.SameSite != "Strict" {
		t.Fatalf("expected canonical same-site to win, got %q", cfg.Security.Auth.Cookie.SameSite)
	}

	if provider.ClientID != "canonical-id" {
		t.Fatalf("expected canonical client-id to win, got %q", provider.ClientID)
	}
	if provider.ClientSecret != "canonical-secret" {
		t.Fatalf("expected canonical client-secret to win, got %q", provider.ClientSecret)
	}
	if provider.AuthURL != "https://canonical.example.com/auth" {
		t.Fatalf("expected canonical auth-url to win, got %q", provider.AuthURL)
	}
	if provider.TokenURL != "https://canonical.example.com/token" {
		t.Fatalf("expected canonical token-url to win, got %q", provider.TokenURL)
	}
	if provider.RedirectURL != "https://canonical.example.com/callback" {
		t.Fatalf("expected canonical redirect-url to win, got %q", provider.RedirectURL)
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
