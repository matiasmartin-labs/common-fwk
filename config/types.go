package config

import "time"

const (
	// JWTAlgorithmHS256 keeps backward-compatible HMAC mode as default.
	JWTAlgorithmHS256 = "HS256"
	// JWTAlgorithmRS256 enables RSA verification mode.
	JWTAlgorithmRS256 = "RS256"

	// RS256KeySourceGenerated derives keys from in-memory generation.
	RS256KeySourceGenerated = "generated"
	// RS256KeySourcePublicPEM uses an explicit PEM-encoded public key.
	RS256KeySourcePublicPEM = "public-pem"
	// RS256KeySourcePrivatePEM uses an explicit PEM-encoded private key.
	RS256KeySourcePrivatePEM = "private-pem"
)

// Config is the root configuration model for common-fwk.
type Config struct {
	Server   ServerConfig
	Security SecurityConfig
	Logging  LoggingConfig
}

// LoggingConfig contains framework logging defaults and per-logger overrides.
type LoggingConfig struct {
	Enabled bool
	Level   string
	Format  string
	Loggers map[string]LoggerOverrideConfig
}

// LoggerOverrideConfig contains per-logger enabled/level overrides.
type LoggerOverrideConfig struct {
	Enabled *bool
	Level   string
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Host           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

// ServerRuntimeLimits groups HTTP runtime limits for server construction.
type ServerRuntimeLimits struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
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
	Algorithm  string
	Secret     string
	Issuer     string
	TTLMinutes int
	RS256      RS256Config
}

// RS256Config contains RSA key resolution settings for RS256 mode.
type RS256Config struct {
	KeySource     string
	KeyID         string
	PublicKeyPEM  string
	PrivateKeyPEM string
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
