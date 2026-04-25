package config

// Config is the root configuration model for common-fwk.
type Config struct {
	Server   ServerConfig
	Security SecurityConfig
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Host string
	Port int
}

// SecurityConfig groups security-related settings.
type SecurityConfig struct {
	Auth AuthConfig
}

// AuthConfig groups authentication and authorization settings.
type AuthConfig struct {
	JWT    JWTConfig
	Cookie CookieConfig
	Login  LoginConfig
	OAuth2 OAuth2Config
}

// JWTConfig contains token signing and lifetime settings.
type JWTConfig struct {
	Secret     string
	Issuer     string
	TTLMinutes int
}

// CookieConfig contains cookie settings used by authentication flows.
type CookieConfig struct {
	Name     string
	Domain   string
	Secure   bool
	HTTPOnly bool
	SameSite string
}

// LoginConfig contains login-specific settings.
type LoginConfig struct {
	Email string
}

// OAuth2Config contains generic OAuth2 provider settings.
type OAuth2Config struct {
	Providers map[string]OAuth2ProviderConfig
}

// OAuth2ProviderConfig stores generic OAuth2 client configuration.
type OAuth2ProviderConfig struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	RedirectURL  string
	Scopes       []string
}
