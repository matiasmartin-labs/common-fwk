package viper

import (
	"errors"
	"strings"

	"github.com/matiasmartin-labs/common-fwk/config"
)

var errEmptyProviderKey = errors.New("provider key must not be empty")

type rawConfig struct {
	Server   rawServerConfig   `mapstructure:"server"`
	Security rawSecurityConfig `mapstructure:"security"`
}

type rawServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
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
	Secret     string `mapstructure:"secret"`
	Issuer     string `mapstructure:"issuer"`
	TTLMinutes int    `mapstructure:"ttlMinutes"`
}

type rawCookieConfig struct {
	Name     string `mapstructure:"name"`
	Domain   string `mapstructure:"domain"`
	Secure   bool   `mapstructure:"secure"`
	HTTPOnly bool   `mapstructure:"httpOnly"`
	SameSite string `mapstructure:"sameSite"`
}

type rawLoginConfig struct {
	Email string `mapstructure:"email"`
}

type rawOAuth2Config struct {
	Providers map[string]rawOAuth2ProviderConfig `mapstructure:"providers"`
}

type rawOAuth2ProviderConfig struct {
	ClientID     string   `mapstructure:"clientID"`
	ClientSecret string   `mapstructure:"clientSecret"`
	AuthURL      string   `mapstructure:"authURL"`
	TokenURL     string   `mapstructure:"tokenURL"`
	RedirectURL  string   `mapstructure:"redirectURL"`
	Scopes       []string `mapstructure:"scopes"`
}

func mapRawToCore(raw rawConfig) (config.Config, error) {
	server := config.NewServerConfig(raw.Server.Host, raw.Server.Port)
	jwt := config.NewJWTConfig(raw.Security.Auth.JWT.Secret, raw.Security.Auth.JWT.Issuer, raw.Security.Auth.JWT.TTLMinutes)
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

	return config.NewConfig(server, security), nil
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
