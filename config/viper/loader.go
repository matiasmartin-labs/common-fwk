package viper

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	setString("security.auth.cookie.name", binding("SECURITY_AUTH_COOKIE_NAME"))
	setString("security.auth.cookie.domain", binding("SECURITY_AUTH_COOKIE_DOMAIN"))
	setString("security.auth.cookie.sameSite", binding("SECURITY_AUTH_COOKIE_SAMESITE"))
	setString("security.auth.login.email", binding("SECURITY_AUTH_LOGIN_EMAIL"))

	if value, ok := snapshot[binding("SERVER_PORT")]; ok {
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse env %q as int: %w", binding("SERVER_PORT"), err)
		}
		v.Set("server.port", parsed)
	}

	if value, ok := snapshot[binding("SECURITY_AUTH_JWT_TTLMINUTES")]; ok {
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse env %q as int: %w", binding("SECURITY_AUTH_JWT_TTLMINUTES"), err)
		}
		v.Set("security.auth.jwt.ttlMinutes", parsed)
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
		v.Set("security.auth.cookie.httpOnly", parsed)
	}

	return nil
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
