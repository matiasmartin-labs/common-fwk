package viper

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/matiasmartin-labs/common-fwk/config"
	"github.com/spf13/viper"
)

// Load reads configuration through the Viper adapter and returns validated core config.
func Load(opts Options) (cfg config.Config, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			cfg = config.Config{}
			err = &LoadError{Err: fmt.Errorf("panic recovered: %v", recovered)}
		}
	}()

	normalized := opts.normalized()
	if normalized.ConfigPath == "" {
		return config.Config{}, &LoadError{Err: fmt.Errorf("config path is required")}
	}

	configType, err := resolveConfigType(normalized.ConfigPath, normalized.ConfigType)
	if err != nil {
		return config.Config{}, &LoadError{Err: err}
	}

	snapshot := environmentSnapshot()
	content, err := os.ReadFile(normalized.ConfigPath)
	if err != nil {
		return config.Config{}, &LoadError{Err: fmt.Errorf("read config file %q: %w", normalized.ConfigPath, err)}
	}

	if normalized.ExpandEnv {
		content = []byte(expandWithSnapshot(string(content), snapshot))
	}

	v := viper.New()
	v.SetConfigType(configType)
	if err := v.ReadConfig(bytes.NewBuffer(content)); err != nil {
		return config.Config{}, &DecodeError{Err: fmt.Errorf("decode %q as %s: %w", normalized.ConfigPath, configType, err)}
	}

	if err := applyEnvironmentOverrides(v, normalized, snapshot); err != nil {
		return config.Config{}, &LoadError{Err: err}
	}

	applyLegacyKeyCompatibility(v)

	var raw rawConfig
	if err := v.Unmarshal(&raw); err != nil {
		return config.Config{}, &DecodeError{Err: fmt.Errorf("unmarshal config: %w", err)}
	}

	if normalized.ExpandEnv {
		raw = expandRawConfig(raw, snapshot)
	}

	mapped, err := mapRawToCore(raw)
	if err != nil {
		return config.Config{}, err
	}

	validated, err := config.ValidateConfig(mapped)
	if err != nil {
		return config.Config{}, &ValidationError{Err: err}
	}

	return validated, nil
}

func applyEnvironmentOverrides(v *viper.Viper, opts Options, snapshot map[string]string) error {
	if !opts.EnvOverride {
		return nil
	}

	binding := func(key string) string {
		prefix := opts.EnvPrefix
		if prefix == "" {
			return key
		}

		return prefix + "_" + key
	}

	setString := func(viperPath, envKey string) {
		if value, ok := snapshot[envKey]; ok {
			v.Set(viperPath, value)
		}
	}

	setString("server.host", binding("SERVER_HOST"))
	setString("security.auth.jwt.secret", binding("SECURITY_AUTH_JWT_SECRET"))
	setString("security.auth.jwt.issuer", binding("SECURITY_AUTH_JWT_ISSUER"))
	setString("security.auth.jwt.algorithm", binding("SECURITY_AUTH_JWT_ALGORITHM"))
	setString("security.auth.jwt.rs256-key-source", binding("SECURITY_AUTH_JWT_RS256_KEY_SOURCE"))
	setString("security.auth.jwt.rs256-key-id", binding("SECURITY_AUTH_JWT_RS256_KEY_ID"))
	setString("security.auth.jwt.rs256-public-key-pem", binding("SECURITY_AUTH_JWT_RS256_PUBLIC_KEY_PEM"))
	setString("security.auth.jwt.rs256-private-key-pem", binding("SECURITY_AUTH_JWT_RS256_PRIVATE_KEY_PEM"))
	setString("security.auth.cookie.name", binding("SECURITY_AUTH_COOKIE_NAME"))
	setString("security.auth.cookie.domain", binding("SECURITY_AUTH_COOKIE_DOMAIN"))
	setString("security.auth.cookie.same-site", binding("SECURITY_AUTH_COOKIE_SAMESITE"))
	setString("security.auth.login.email", binding("SECURITY_AUTH_LOGIN_EMAIL"))

	if value, ok := snapshot[binding("SERVER_PORT")]; ok {
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse env %q as int: %w", binding("SERVER_PORT"), err)
		}
		v.Set("server.port", parsed)
	}

	if value, ok := snapshot[binding("SERVER_READ_TIMEOUT")]; ok {
		parsed, err := time.ParseDuration(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse env %q as duration: %w", binding("SERVER_READ_TIMEOUT"), err)
		}
		v.Set("server.read-timeout", parsed)
	}

	if value, ok := snapshot[binding("SERVER_WRITE_TIMEOUT")]; ok {
		parsed, err := time.ParseDuration(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse env %q as duration: %w", binding("SERVER_WRITE_TIMEOUT"), err)
		}
		v.Set("server.write-timeout", parsed)
	}

	if value, ok := snapshot[binding("SERVER_MAX_HEADER_BYTES")]; ok {
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse env %q as int: %w", binding("SERVER_MAX_HEADER_BYTES"), err)
		}
		v.Set("server.max-header-bytes", parsed)
	}

	if value, ok := snapshot[binding("SECURITY_AUTH_JWT_TTLMINUTES")]; ok {
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse env %q as int: %w", binding("SECURITY_AUTH_JWT_TTLMINUTES"), err)
		}
		v.Set("security.auth.jwt.ttl-minutes", parsed)
	}

	if value, ok := snapshot[binding("SECURITY_AUTH_COOKIE_SECURE")]; ok {
		parsed, err := strconv.ParseBool(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse env %q as bool: %w", binding("SECURITY_AUTH_COOKIE_SECURE"), err)
		}
		v.Set("security.auth.cookie.secure", parsed)
	}

	if value, ok := snapshot[binding("SECURITY_AUTH_COOKIE_HTTPONLY")]; ok {
		parsed, err := strconv.ParseBool(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse env %q as bool: %w", binding("SECURITY_AUTH_COOKIE_HTTPONLY"), err)
		}
		v.Set("security.auth.cookie.http-only", parsed)
	}

	return nil
}

func applyLegacyKeyCompatibility(v *viper.Viper) {
	aliases := [][2]string{
		{"security.auth.jwt.ttl-minutes", "security.auth.jwt.ttlMinutes"},
		{"security.auth.jwt.rs256-key-source", "security.auth.jwt.rs256KeySource"},
		{"security.auth.jwt.rs256-key-id", "security.auth.jwt.rs256KeyID"},
		{"security.auth.jwt.rs256-public-key-pem", "security.auth.jwt.rs256PublicKeyPEM"},
		{"security.auth.jwt.rs256-private-key-pem", "security.auth.jwt.rs256PrivateKeyPEM"},
		{"security.auth.cookie.http-only", "security.auth.cookie.httpOnly"},
		{"security.auth.cookie.same-site", "security.auth.cookie.sameSite"},
	}

	for _, alias := range aliases {
		canonical := alias[0]
		legacy := alias[1]
		if !v.IsSet(canonical) && v.IsSet(legacy) {
			v.Set(canonical, v.Get(legacy))
		}
	}

	providers, ok := providerSettings(v)
	if !ok {
		return
	}

	for providerKey := range providers {
		applyProviderAlias(v, providerKey, "client-id", "clientID")
		applyProviderAlias(v, providerKey, "client-secret", "clientSecret")
		applyProviderAlias(v, providerKey, "auth-url", "authURL")
		applyProviderAlias(v, providerKey, "token-url", "tokenURL")
		applyProviderAlias(v, providerKey, "redirect-url", "redirectURL")
	}
}

func providerSettings(v *viper.Viper) (map[string]any, bool) {
	securityRaw, ok := v.AllSettings()["security"]
	if !ok {
		return nil, false
	}

	security, ok := securityRaw.(map[string]any)
	if !ok {
		return nil, false
	}

	authRaw, ok := security["auth"]
	if !ok {
		return nil, false
	}

	auth, ok := authRaw.(map[string]any)
	if !ok {
		return nil, false
	}

	oauth2Raw, ok := auth["oauth2"]
	if !ok {
		return nil, false
	}

	oauth2, ok := oauth2Raw.(map[string]any)
	if !ok {
		return nil, false
	}

	providersRaw, ok := oauth2["providers"]
	if !ok {
		return nil, false
	}

	providers, ok := providersRaw.(map[string]any)
	if !ok {
		return nil, false
	}

	return providers, true
}

func applyProviderAlias(v *viper.Viper, providerKey, canonicalKey, legacyKey string) {
	canonicalPath := fmt.Sprintf("security.auth.oauth2.providers.%s.%s", providerKey, canonicalKey)
	legacyPath := fmt.Sprintf("security.auth.oauth2.providers.%s.%s", providerKey, legacyKey)
	if !v.IsSet(canonicalPath) && v.IsSet(legacyPath) {
		v.Set(canonicalPath, v.Get(legacyPath))
	}
}

func environmentSnapshot() map[string]string {
	env := os.Environ()
	snapshot := make(map[string]string, len(env))
	for _, pair := range env {
		idx := strings.Index(pair, "=")
		if idx < 0 {
			continue
		}

		snapshot[pair[:idx]] = pair[idx+1:]
	}

	return snapshot
}

func expandWithSnapshot(input string, snapshot map[string]string) string {
	return os.Expand(input, func(name string) string {
		if value, ok := snapshot[name]; ok {
			return value
		}

		return ""
	})
}

func expandRawConfig(raw rawConfig, snapshot map[string]string) rawConfig {
	raw.Server.Host = expandWithSnapshot(raw.Server.Host, snapshot)
	raw.Security.Auth.JWT.Secret = expandWithSnapshot(raw.Security.Auth.JWT.Secret, snapshot)
	raw.Security.Auth.JWT.Issuer = expandWithSnapshot(raw.Security.Auth.JWT.Issuer, snapshot)
	raw.Security.Auth.JWT.Algorithm = expandWithSnapshot(raw.Security.Auth.JWT.Algorithm, snapshot)
	raw.Security.Auth.JWT.RS256KeySource = expandWithSnapshot(raw.Security.Auth.JWT.RS256KeySource, snapshot)
	raw.Security.Auth.JWT.RS256KeyID = expandWithSnapshot(raw.Security.Auth.JWT.RS256KeyID, snapshot)
	raw.Security.Auth.JWT.RS256PublicKeyPEM = expandWithSnapshot(raw.Security.Auth.JWT.RS256PublicKeyPEM, snapshot)
	raw.Security.Auth.JWT.RS256PrivateKeyPEM = expandWithSnapshot(raw.Security.Auth.JWT.RS256PrivateKeyPEM, snapshot)
	raw.Security.Auth.Cookie.Name = expandWithSnapshot(raw.Security.Auth.Cookie.Name, snapshot)
	raw.Security.Auth.Cookie.Domain = expandWithSnapshot(raw.Security.Auth.Cookie.Domain, snapshot)
	raw.Security.Auth.Cookie.SameSite = expandWithSnapshot(raw.Security.Auth.Cookie.SameSite, snapshot)
	raw.Security.Auth.Login.Email = expandWithSnapshot(raw.Security.Auth.Login.Email, snapshot)

	expandedProviders := make(map[string]rawOAuth2ProviderConfig, len(raw.Security.Auth.OAuth2.Providers))
	for key, provider := range raw.Security.Auth.OAuth2.Providers {
		expandedScopes := make([]string, len(provider.Scopes))
		for index, scope := range provider.Scopes {
			expandedScopes[index] = expandWithSnapshot(scope, snapshot)
		}

		expandedProviders[key] = rawOAuth2ProviderConfig{
			ClientID:     expandWithSnapshot(provider.ClientID, snapshot),
			ClientSecret: expandWithSnapshot(provider.ClientSecret, snapshot),
			AuthURL:      expandWithSnapshot(provider.AuthURL, snapshot),
			TokenURL:     expandWithSnapshot(provider.TokenURL, snapshot),
			RedirectURL:  expandWithSnapshot(provider.RedirectURL, snapshot),
			Scopes:       expandedScopes,
		}
	}
	raw.Security.Auth.OAuth2.Providers = expandedProviders

	return raw
}
