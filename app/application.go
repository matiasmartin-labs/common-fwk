package app

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matiasmartin-labs/common-fwk/config"
	httpgin "github.com/matiasmartin-labs/common-fwk/http/gin"
	"github.com/matiasmartin-labs/common-fwk/logging"
	loggingslog "github.com/matiasmartin-labs/common-fwk/logging/slog"
	"github.com/matiasmartin-labs/common-fwk/security"
	securityjwt "github.com/matiasmartin-labs/common-fwk/security/jwt"
)

var (
	ErrServerNotReady       = errors.New("server not ready")
	ErrSecurityNotReady     = errors.New("security not ready")
	ErrInvalidPath          = errors.New("invalid path")
	ErrNilHandler           = errors.New("nil handler")
	ErrNilListener          = errors.New("nil listener")
	ErrRouteConflict        = errors.New("route conflict")
	ErrInvalidPresetOptions = errors.New("invalid preset options")
	ErrLoggingNotReady      = errors.New("logging runtime not ready")
	ErrLoggerNameRequired   = errors.New("logger name is required")
)

const (
	defaultHealthPath = "/healthz"
	defaultReadyPath  = "/readyz"
)

// ReadinessCheck is a synchronous readiness probe.
//
// Returning nil means the check passed. Returning an error marks the
// application as not-ready for the current request.
type ReadinessCheck func() error

// HealthReadinessOptions configures health/readiness preset registration.
type HealthReadinessOptions struct {
	// HealthPath overrides the health endpoint path.
	//
	// Defaults to "/healthz" when omitted.
	HealthPath string
	// ReadyPath overrides the readiness endpoint path.
	//
	// Defaults to "/readyz" when omitted.
	ReadyPath string
	// Checks are evaluated synchronously in order for each readiness request.
	Checks []ReadinessCheck
}

// Application is an instance-scoped bootstrap container.
type Application struct {
	cfg            config.Config
	server         http.Server
	handler        *gin.Engine
	validator      security.Validator
	loggerRegistry logging.Registry
	logOutput      io.Writer
	serverReady    bool
	securityReady  bool
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
		handler:   h,
		logOutput: io.Discard,
	}
}

// UseConfig sets application configuration and returns the same instance.
func (a *Application) UseConfig(cfg config.Config) *Application {
	a.cfg = cfg
	a.wireLoggingRuntime(cfg.Logging)
	return a
}

func (a *Application) wireLoggingRuntime(loggingCfg config.LoggingConfig) {
	a.loggerRegistry = loggingslog.NewRegistry(normalizeLoggingConfig(loggingCfg), a.logOutput)
}

func normalizeLoggingConfig(cfg config.LoggingConfig) config.LoggingConfig {
	if strings.TrimSpace(cfg.Level) == "" {
		cfg.Level = "info"
	}
	if strings.TrimSpace(cfg.Format) == "" {
		cfg.Format = "json"
	}
	if cfg.Loggers == nil {
		cfg.Loggers = map[string]config.LoggerOverrideConfig{}
	}

	return cfg
}

// GetLogger returns a deterministic named logger when runtime is ready.
func (a *Application) GetLogger(name string) (logging.Logger, error) {
	if a.loggerRegistry == nil {
		return nil, fmt.Errorf("get logger %q: %w", name, ErrLoggingNotReady)
	}

	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, fmt.Errorf("get logger %q: %w", name, ErrLoggerNameRequired)
	}

	return a.loggerRegistry.Get(trimmed), nil
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

func resolveHealthReadinessOptions(opts HealthReadinessOptions) (HealthReadinessOptions, error) {
	resolved := HealthReadinessOptions{
		HealthPath: defaultHealthPath,
		ReadyPath:  defaultReadyPath,
		Checks:     opts.Checks,
	}

	if opts.HealthPath != "" {
		if strings.TrimSpace(opts.HealthPath) == "" {
			return HealthReadinessOptions{}, fmt.Errorf("health path is blank: %w", ErrInvalidPresetOptions)
		}
		resolved.HealthPath = opts.HealthPath
	}

	if opts.ReadyPath != "" {
		if strings.TrimSpace(opts.ReadyPath) == "" {
			return HealthReadinessOptions{}, fmt.Errorf("ready path is blank: %w", ErrInvalidPresetOptions)
		}
		resolved.ReadyPath = opts.ReadyPath
	}

	if resolved.HealthPath == resolved.ReadyPath {
		return HealthReadinessOptions{}, fmt.Errorf(
			"health and readiness paths must differ: %w",
			ErrInvalidPresetOptions,
		)
	}

	return resolved, nil
}

func hasGETRoute(routes gin.RoutesInfo, path string) bool {
	for _, route := range routes {
		if route.Method == http.MethodGet && route.Path == path {
			return true
		}
	}

	return false
}

func (a *Application) ensureNoPresetRouteConflict(paths ...string) error {
	routes := a.handler.Routes()
	for _, path := range paths {
		if hasGETRoute(routes, path) {
			return fmt.Errorf("method=%s path=%q: %w", http.MethodGet, path, ErrRouteConflict)
		}
	}

	return nil
}

func (a *Application) readinessInvariantSatisfied() bool {
	return a.serverReady && a.handler != nil && a.server.Handler != nil
}

func (a *Application) readinessStatus(checks []ReadinessCheck) int {
	if !a.readinessInvariantSatisfied() {
		return http.StatusServiceUnavailable
	}

	for _, check := range checks {
		if check == nil {
			return http.StatusServiceUnavailable
		}

		if err := check(); err != nil {
			return http.StatusServiceUnavailable
		}
	}

	return http.StatusOK
}

// EnableHealthReadinessPresets registers health/readiness handlers explicitly.
//
// The method is opt-in and never runs implicitly from UseServer. It requires a
// ready server bootstrap and fails when requested paths are invalid or conflict
// with already registered GET routes.
func (a *Application) EnableHealthReadinessPresets(opts HealthReadinessOptions) error {
	if err := a.ensureServerReady(); err != nil {
		return err
	}

	resolvedOpts, err := resolveHealthReadinessOptions(opts)
	if err != nil {
		return err
	}

	if err := a.ensureNoPresetRouteConflict(resolvedOpts.HealthPath, resolvedOpts.ReadyPath); err != nil {
		return err
	}

	a.handler.GET(resolvedOpts.HealthPath, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	a.handler.GET(resolvedOpts.ReadyPath, func(c *gin.Context) {
		c.Status(a.readinessStatus(resolvedOpts.Checks))
	})

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
