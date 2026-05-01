package app

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matiasmartin-labs/common-fwk/config"
	httpgin "github.com/matiasmartin-labs/common-fwk/http/gin"
	"github.com/matiasmartin-labs/common-fwk/security"
	securityjwt "github.com/matiasmartin-labs/common-fwk/security/jwt"
)

var (
	ErrServerNotReady   = errors.New("server not ready")
	ErrSecurityNotReady = errors.New("security not ready")
	ErrInvalidPath      = errors.New("invalid path")
	ErrNilHandler       = errors.New("nil handler")
	ErrNilListener      = errors.New("nil listener")
)

// Application is an instance-scoped bootstrap container.
type Application struct {
	cfg           config.Config
	server        http.Server
	handler       *gin.Engine
	validator     security.Validator
	serverReady   bool
	securityReady bool
}

// GetConfig returns a read-only snapshot of the current runtime config.
func (a *Application) GetConfig() config.Config {
	return cloneConfig(a.cfg)
}

// GetSecurityValidator returns the current security validator when wired.
func (a *Application) GetSecurityValidator() security.Validator {
	return a.validator
}

// IsSecurityReady reports whether security runtime was fully wired.
func (a *Application) IsSecurityReady() bool {
	return a.securityReady
}

// NewApplication creates a bootstrap instance with safe defaults.
func NewApplication() *Application {
	h := gin.New()

	return &Application{
		handler: h,
	}
}

// UseConfig sets application configuration and returns the same instance.
func (a *Application) UseConfig(cfg config.Config) *Application {
	a.cfg = cfg
	return a
}

func cloneConfig(cfg config.Config) config.Config {
	cfg.Security.Auth.OAuth2.Providers = cloneOAuth2Providers(cfg.Security.Auth.OAuth2.Providers)

	return cfg
}

func cloneOAuth2Providers(providers map[string]config.OAuth2ProviderConfig) map[string]config.OAuth2ProviderConfig {
	if providers == nil {
		return nil
	}

	cloned := make(map[string]config.OAuth2ProviderConfig, len(providers))
	for key, provider := range providers {
		provider.Scopes = cloneStringSlice(provider.Scopes)
		cloned[key] = provider
	}

	return cloned
}

func cloneStringSlice(values []string) []string {
	if values == nil {
		return nil
	}

	cloned := make([]string, len(values))
	copy(cloned, values)

	return cloned
}

// UseServer wires the HTTP server and marks server readiness.
func (a *Application) UseServer() *Application {
	if a.handler == nil {
		a.handler = gin.New()
	}

	a.server.Handler = a.handler
	a.server.ReadTimeout = a.cfg.Server.ReadTimeout
	a.server.WriteTimeout = a.cfg.Server.WriteTimeout
	a.server.MaxHeaderBytes = a.cfg.Server.MaxHeaderBytes
	a.serverReady = true

	return a
}

// UseServerSecurity sets the validator and marks security readiness.
func (a *Application) UseServerSecurity(v security.Validator) *Application {
	a.validator = v
	a.securityReady = v != nil

	return a
}

// UseServerSecurityFromConfig wires validator from currently loaded config.
func (a *Application) UseServerSecurityFromConfig() (*Application, error) {
	validatedCfg, err := config.ValidateConfig(a.cfg)
	if err != nil {
		return a, fmt.Errorf("validate config before security wiring: %w", err)
	}

	compat, err := securityjwt.FromConfigJWT(validatedCfg.Security.Auth.JWT)
	if err != nil {
		return a, fmt.Errorf("build validator options from config: %w", err)
	}

	validator, err := securityjwt.NewValidator(compat.Options)
	if err != nil {
		return a, fmt.Errorf("build jwt validator: %w", err)
	}

	a.UseServerSecurity(validator)
	return a, nil
}

func (a *Application) ensureServerReady() error {
	if !a.serverReady {
		return ErrServerNotReady
	}

	return nil
}

func (a *Application) ensureSecurityReady() error {
	if !a.securityReady {
		return ErrSecurityNotReady
	}

	return nil
}

func validatePath(path string) error {
	if strings.TrimSpace(path) == "" {
		return ErrInvalidPath
	}

	return nil
}

func validateHandler(h gin.HandlerFunc) error {
	if h == nil {
		return ErrNilHandler
	}

	return nil
}

// RegisterGET registers a GET route in the application's gin engine.
func (a *Application) RegisterGET(path string, h gin.HandlerFunc) error {
	if err := a.ensureServerReady(); err != nil {
		return err
	}
	if err := validatePath(path); err != nil {
		return err
	}
	if err := validateHandler(h); err != nil {
		return err
	}

	a.handler.GET(path, h)

	return nil
}

// RegisterPOST registers a POST route in the application's gin engine.
func (a *Application) RegisterPOST(path string, h gin.HandlerFunc) error {
	if err := a.ensureServerReady(); err != nil {
		return err
	}
	if err := validatePath(path); err != nil {
		return err
	}
	if err := validateHandler(h); err != nil {
		return err
	}

	a.handler.POST(path, h)

	return nil
}

// RegisterProtectedGET registers a GET route guarded by auth middleware.
func (a *Application) RegisterProtectedGET(path string, h gin.HandlerFunc) error {
	if err := a.ensureServerReady(); err != nil {
		return err
	}
	if err := a.ensureSecurityReady(); err != nil {
		return err
	}
	if err := validatePath(path); err != nil {
		return err
	}
	if err := validateHandler(h); err != nil {
		return err
	}

	a.handler.GET(path, httpgin.NewAuthMiddleware(a.validator), h)

	return nil
}

// Run starts serving with ListenAndServe.
func (a *Application) Run() error {
	if err := a.ensureServerReady(); err != nil {
		return err
	}

	a.server.Handler = a.handler
	a.server.Addr = net.JoinHostPort(a.cfg.Server.Host, strconv.Itoa(a.cfg.Server.Port))

	return a.server.ListenAndServe()
}

// RunListener starts serving on a provided listener.
func (a *Application) RunListener(l net.Listener) error {
	if err := a.ensureServerReady(); err != nil {
		return err
	}
	if l == nil {
		return ErrNilListener
	}

	a.server.Handler = a.handler

	return a.server.Serve(l)
}
