package config

import "time"

const (
	defaultServerHost = "127.0.0.1"
	defaultServerPort = 8080

	defaultServerReadTimeout    = 10 * time.Second
	defaultServerWriteTimeout   = 10 * time.Second
	defaultServerMaxHeaderBytes = 1 << 20

	defaultJWTTTLMinutes = 60

	defaultCookieName     = "session"
	defaultCookieSameSite = "Lax"
)

// NewConfig constructs the root config from explicit dependencies.
func NewConfig(server ServerConfig, security SecurityConfig) Config {
	return Config{
		Server:   server,
		Security: security,
	}
}

// NewServerConfig returns a server config with useful defaults.
func NewServerConfig(host string, port int, limits ...ServerRuntimeLimits) ServerConfig {
	runtime := ServerRuntimeLimits{}
	if len(limits) > 0 {
		runtime = limits[0]
	}

	if host == "" {
		host = defaultServerHost
	}

	if port == 0 {
		port = defaultServerPort
	}

	if runtime.ReadTimeout == 0 {
		runtime.ReadTimeout = defaultServerReadTimeout
	}

	if runtime.WriteTimeout == 0 {
		runtime.WriteTimeout = defaultServerWriteTimeout
	}

	if runtime.MaxHeaderBytes == 0 {
		runtime.MaxHeaderBytes = defaultServerMaxHeaderBytes
	}

	return ServerConfig{
		Host:           host,
		Port:           port,
		ReadTimeout:    runtime.ReadTimeout,
		WriteTimeout:   runtime.WriteTimeout,
		MaxHeaderBytes: runtime.MaxHeaderBytes,
	}
}

// NewSecurityConfig returns security config from explicit auth config input.
func NewSecurityConfig(auth AuthConfig) SecurityConfig {
	return SecurityConfig{Auth: auth}
}

// NewAuthConfig returns auth config from explicit nested config inputs.
func NewAuthConfig(jwt JWTConfig, cookie CookieConfig, login LoginConfig, oauth2 OAuth2Config) AuthConfig {
	return AuthConfig{
		JWT:    jwt,
		Cookie: cookie,
		Login:  login,
		OAuth2: oauth2,
	}
}

// NewJWTConfig returns JWT config with useful defaults.
func NewJWTConfig(secret, issuer string, ttlMinutes int) JWTConfig {
	if ttlMinutes == 0 {
		ttlMinutes = defaultJWTTTLMinutes
	}

	return JWTConfig{
		Secret:     secret,
		Issuer:     issuer,
		TTLMinutes: ttlMinutes,
	}
}

// NewCookieConfig returns cookie config with useful defaults.
func NewCookieConfig(name, domain string, secure, httpOnly bool, sameSite string) CookieConfig {
	if name == "" {
		name = defaultCookieName
	}

	if sameSite == "" {
		sameSite = defaultCookieSameSite
	}

	return CookieConfig{
		Name:     name,
		Domain:   domain,
		Secure:   secure,
		HTTPOnly: httpOnly,
		SameSite: sameSite,
	}
}

// NewLoginConfig returns login config with explicit email input.
func NewLoginConfig(email string) LoginConfig {
	return LoginConfig{Email: email}
}

// NewOAuth2Config returns OAuth2 config with a defensive copy of providers.
func NewOAuth2Config(providers map[string]OAuth2ProviderConfig) OAuth2Config {
	clonedProviders := make(map[string]OAuth2ProviderConfig, len(providers))
	for key, provider := range providers {
		clonedProviders[key] = NewOAuth2ProviderConfig(
			provider.ClientID,
			provider.ClientSecret,
			provider.AuthURL,
			provider.TokenURL,
			provider.RedirectURL,
			provider.Scopes,
		)
	}

	return OAuth2Config{Providers: clonedProviders}
}

// NewOAuth2ProviderConfig returns provider config with a defensive scopes copy.
func NewOAuth2ProviderConfig(clientID, clientSecret, authURL, tokenURL, redirectURL string, scopes []string) OAuth2ProviderConfig {
	clonedScopes := make([]string, len(scopes))
	copy(clonedScopes, scopes)

	return OAuth2ProviderConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		RedirectURL:  redirectURL,
		Scopes:       clonedScopes,
	}
}
