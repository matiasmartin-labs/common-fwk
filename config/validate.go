package config

import (
	"fmt"
	"net/mail"
	"strings"
)

var allowedSameSiteValues = map[string]struct{}{
	"Lax":    {},
	"Strict": {},
	"None":   {},
}

// ValidateConfig validates a configuration value and returns a normalized copy.
func ValidateConfig(cfg Config) (Config, error) {
	normalized := cfg
	normalized.Security.Auth.Login.Email = normalizeLoginEmail(normalized.Security.Auth.Login.Email)

	if err := validateServer(normalized.Server); err != nil {
		return Config{}, wrapInvalidConfig(err)
	}

	if err := validateJWT(normalized.Security.Auth.JWT); err != nil {
		return Config{}, wrapInvalidConfig(err)
	}

	if err := validateCookie(normalized.Security.Auth.Cookie); err != nil {
		return Config{}, wrapInvalidConfig(err)
	}

	if err := validateLogin(normalized.Security.Auth.Login); err != nil {
		return Config{}, wrapInvalidConfig(err)
	}

	if err := validateOAuth2(normalized.Security.Auth.OAuth2); err != nil {
		return Config{}, wrapInvalidConfig(err)
	}

	return normalized, nil
}

func validateServer(cfg ServerConfig) error {
	if strings.TrimSpace(cfg.Host) == "" {
		return invalidAt("server.host", ErrRequired)
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return invalidAt("server.port", fmt.Errorf("%w: must be between 1 and 65535", ErrOutOfRange))
	}

	return nil
}

func validateJWT(cfg JWTConfig) error {
	if strings.TrimSpace(cfg.Secret) == "" {
		return invalidAt("security.auth.jwt.secret", ErrRequired)
	}

	if strings.TrimSpace(cfg.Issuer) == "" {
		return invalidAt("security.auth.jwt.issuer", ErrRequired)
	}

	if cfg.TTLMinutes < 1 {
		return invalidAt("security.auth.jwt.ttlMinutes", fmt.Errorf("%w: must be positive", ErrOutOfRange))
	}

	return nil
}

func validateCookie(cfg CookieConfig) error {
	if strings.TrimSpace(cfg.Name) == "" {
		return invalidAt("security.auth.cookie.name", ErrRequired)
	}

	if cfg.SameSite == "" {
		return invalidAt("security.auth.cookie.sameSite", ErrRequired)
	}

	if _, ok := allowedSameSiteValues[cfg.SameSite]; !ok {
		return invalidAt("security.auth.cookie.sameSite", fmt.Errorf("%w: got %q", ErrOutOfRange, cfg.SameSite))
	}

	return nil
}

func validateLogin(cfg LoginConfig) error {
	if cfg.Email == "" {
		return invalidAt("security.auth.login.email", ErrRequired)
	}

	if _, err := mail.ParseAddress(cfg.Email); err != nil {
		return invalidAt("security.auth.login.email", fmt.Errorf("%w: %v", ErrInvalidEmail, err))
	}

	return nil
}

func validateOAuth2(cfg OAuth2Config) error {
	for providerKey, provider := range cfg.Providers {
		basePath := fmt.Sprintf("security.auth.oauth2.providers.%s", providerKey)
		if strings.TrimSpace(providerKey) == "" {
			return invalidAt("security.auth.oauth2.providers", fmt.Errorf("%w: provider key must not be empty", ErrRequired))
		}

		if strings.TrimSpace(provider.ClientID) == "" {
			return invalidAt(basePath+".clientID", ErrRequired)
		}

		if strings.TrimSpace(provider.ClientSecret) == "" {
			return invalidAt(basePath+".clientSecret", ErrRequired)
		}

		if strings.TrimSpace(provider.AuthURL) == "" {
			return invalidAt(basePath+".authURL", ErrRequired)
		}

		if strings.TrimSpace(provider.TokenURL) == "" {
			return invalidAt(basePath+".tokenURL", ErrRequired)
		}

		if strings.TrimSpace(provider.RedirectURL) == "" {
			return invalidAt(basePath+".redirectURL", ErrRequired)
		}
	}

	return nil
}

func normalizeLoginEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
