package viper

import (
	"errors"
	"strings"
	"time"

	"github.com/matiasmartin-labs/common-fwk/config"
)

var errEmptyProviderKey = errors.New("provider key must not be empty")

type rawConfig struct {
	Server   rawServerConfig   `mapstructure:"server"`
	Security rawSecurityConfig `mapstructure:"security"`
	Logging  rawLoggingConfig  `mapstructure:"logging"`
}

type rawServerConfig struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"read-timeout"`
	WriteTimeout   time.Duration `mapstructure:"write-timeout"`
	MaxHeaderBytes int           `mapstructure:"max-header-bytes"`
}

type rawSecurityConfig struct {
	Auth rawAuthConfig `mapstructure:"auth"`
}

type rawAuthConfig struct {
	JWT    rawJWTConfig                    `mapstructure:"jwt"`
	Cookie rawCookieConfig                 `mapstructure:"cookie"`
	Login  rawLoginConfig                  `mapstructure:"login"`
	OAuth2 rawOAuth2Config                 `mapstructure:"oauth2"`
	Legacy map[string]map[string]rawLegacy `mapstructure:",remain"`
}

// rawLegacy captures unrecognized fields without impacting mapping behavior.
type rawLegacy struct{}

type rawJWTConfig struct {
	Algorithm          string `mapstructure:"algorithm"`
	Secret             string `mapstructure:"secret"`
	Issuer             string `mapstructure:"issuer"`
	TTLMinutes         int    `mapstructure:"ttl-minutes"`
	RS256KeySource     string `mapstructure:"rs256-key-source"`
	RS256KeyID         string `mapstructure:"rs256-key-id"`
	RS256PublicKeyPEM  string `mapstructure:"rs256-public-key-pem"`
	RS256PrivateKeyPEM string `mapstructure:"rs256-private-key-pem"`
}

type rawCookieConfig struct {
	Name     string `mapstructure:"name"`
	Domain   string `mapstructure:"domain"`
	Secure   bool   `mapstructure:"secure"`
	HTTPOnly bool   `mapstructure:"http-only"`
	SameSite string `mapstructure:"same-site"`
}

type rawLoginConfig struct {
	Email string `mapstructure:"email"`
}

type rawOAuth2Config struct {
	Providers map[string]rawOAuth2ProviderConfig `mapstructure:"providers"`
}

type rawOAuth2ProviderConfig struct {
	ClientID     string   `mapstructure:"client-id"`
	ClientSecret string   `mapstructure:"client-secret"`
	AuthURL      string   `mapstructure:"auth-url"`
	TokenURL     string   `mapstructure:"token-url"`
	RedirectURL  string   `mapstructure:"redirect-url"`
	Scopes       []string `mapstructure:"scopes"`
}

type rawLoggingConfig struct {
	Enabled *bool                              `mapstructure:"enabled"`
	Level   string                             `mapstructure:"level"`
	Format  string                             `mapstructure:"format"`
	Loggers map[string]rawLoggerOverrideConfig `mapstructure:"loggers"`
}

type rawLoggerOverrideConfig struct {
	Enabled *bool  `mapstructure:"enabled"`
	Level   string `mapstructure:"level"`
}

func mapRawToCore(raw rawConfig) (config.Config, error) {
	server := config.NewServerConfig(
		raw.Server.Host,
		raw.Server.Port,
		config.ServerRuntimeLimits{
			ReadTimeout:    raw.Server.ReadTimeout,
			WriteTimeout:   raw.Server.WriteTimeout,
			MaxHeaderBytes: raw.Server.MaxHeaderBytes,
		},
	)
	jwt := config.NewJWTConfig(raw.Security.Auth.JWT.Secret, raw.Security.Auth.JWT.Issuer, raw.Security.Auth.JWT.TTLMinutes)
	if raw.Security.Auth.JWT.Algorithm != "" {
		jwt.Algorithm = raw.Security.Auth.JWT.Algorithm
	}
	jwt.RS256 = config.RS256Config{
		KeySource:     raw.Security.Auth.JWT.RS256KeySource,
		KeyID:         raw.Security.Auth.JWT.RS256KeyID,
		PublicKeyPEM:  raw.Security.Auth.JWT.RS256PublicKeyPEM,
		PrivateKeyPEM: raw.Security.Auth.JWT.RS256PrivateKeyPEM,
	}
	cookie := config.NewCookieConfig(
		raw.Security.Auth.Cookie.Name,
		raw.Security.Auth.Cookie.Domain,
		raw.Security.Auth.Cookie.Secure,
		raw.Security.Auth.Cookie.HTTPOnly,
		raw.Security.Auth.Cookie.SameSite,
	)
	login := config.NewLoginConfig(raw.Security.Auth.Login.Email)

	oauth2Providers, err := mapProviders(raw.Security.Auth.OAuth2.Providers)
	if err != nil {
		return config.Config{}, err
	}

	oauth2 := config.NewOAuth2Config(oauth2Providers)
	auth := config.NewAuthConfig(jwt, cookie, login, oauth2)
	security := config.NewSecurityConfig(auth)

	loggingEnabled := true
	if raw.Logging.Enabled != nil {
		loggingEnabled = *raw.Logging.Enabled
	}

	loggerOverrides := make(map[string]config.LoggerOverrideConfig, len(raw.Logging.Loggers))
	for loggerName, loggerOverride := range raw.Logging.Loggers {
		loggerOverrides[loggerName] = config.LoggerOverrideConfig{
			Enabled: loggerOverride.Enabled,
			Level:   loggerOverride.Level,
		}
	}

	loggingCfg := config.NewLoggingConfig(loggingEnabled, raw.Logging.Level, raw.Logging.Format, loggerOverrides)

	return config.NewConfig(server, security, loggingCfg), nil
}

func mapProviders(providers map[string]rawOAuth2ProviderConfig) (map[string]config.OAuth2ProviderConfig, error) {
	if len(providers) == 0 {
		return map[string]config.OAuth2ProviderConfig{}, nil
	}

	mapped := make(map[string]config.OAuth2ProviderConfig, len(providers))
	for rawKey, provider := range providers {
		key := strings.TrimSpace(rawKey)
		if key == "" {
			return nil, &MappingError{Path: "security.auth.oauth2.providers", Err: errEmptyProviderKey}
		}

		scopes := make([]string, len(provider.Scopes))
		copy(scopes, provider.Scopes)

		mapped[key] = config.NewOAuth2ProviderConfig(
			provider.ClientID,
			provider.ClientSecret,
			provider.AuthURL,
			provider.TokenURL,
			provider.RedirectURL,
			scopes,
		)
	}

	return mapped, nil
}
