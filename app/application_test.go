package app

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matiasmartin-labs/common-fwk/config"
	"github.com/matiasmartin-labs/common-fwk/security/claims"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type fakeValidator struct {
	validateFn func(context.Context, string) (claims.Claims, error)
}

func (f *fakeValidator) Validate(ctx context.Context, raw string) (claims.Claims, error) {
	if f.validateFn != nil {
		return f.validateFn(ctx, raw)
	}

	return claims.Claims{}, nil
}

func testConfig() config.Config {
	return config.Config{
		Server: config.ServerConfig{Host: "127.0.0.1", Port: 0, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, MaxHeaderBytes: 1 << 20},
	}
}

func testConfigWithOAuth2Provider() config.Config {
	return config.NewConfig(
		config.NewServerConfig("127.0.0.1", 8080),
		config.NewSecurityConfig(config.NewAuthConfig(
			config.NewJWTConfig("secret", "common-fwk", 15),
			config.NewCookieConfig("session", "example.com", true, true, "Lax"),
			config.NewLoginConfig("owner@example.com"),
			config.NewOAuth2Config(map[string]config.OAuth2ProviderConfig{
				"github": config.NewOAuth2ProviderConfig(
					"client-id",
					"client-secret",
					"https://github.com/login/oauth/authorize",
					"https://github.com/login/oauth/access_token",
					"https://app.example.com/auth/github/callback",
					[]string{"read:user", "user:email"},
				),
			}),
		)),
	)
}

func TestBootstrapChain_PreservesPointerAndReadiness(t *testing.T) {
	a := NewApplication()
	original := a

	validator := &fakeValidator{}
	got := a.UseConfig(testConfig()).UseServer().UseServerSecurity(validator)

	if got != original {
		t.Fatalf("expected fluent chain to preserve same pointer")
	}
	if !a.serverReady {
		t.Fatalf("expected serverReady=true")
	}
	if !a.securityReady {
		t.Fatalf("expected securityReady=true")
	}
	if a.handler == nil {
		t.Fatalf("expected handler to be initialized")
	}
	if a.server.Handler == nil {
		t.Fatalf("expected server handler wiring")
	}
	if a.server.ReadTimeout != 10*time.Second {
		t.Fatalf("expected server read timeout to be wired, got %s", a.server.ReadTimeout)
	}
	if a.server.WriteTimeout != 10*time.Second {
		t.Fatalf("expected server write timeout to be wired, got %s", a.server.WriteTimeout)
	}
	if a.server.MaxHeaderBytes != 1<<20 {
		t.Fatalf("expected server max header bytes to be wired, got %d", a.server.MaxHeaderBytes)
	}
}

func TestUseServer_WiresRuntimeLimitsFromConfig(t *testing.T) {
	tests := []struct {
		name             string
		serverConfig     config.ServerConfig
		wantReadTimeout  time.Duration
		wantWriteTimeout time.Duration
		wantHeaderBytes  int
	}{
		{
			name: "explicit values",
			serverConfig: config.ServerConfig{
				Host:           "127.0.0.1",
				Port:           8080,
				ReadTimeout:    2 * time.Second,
				WriteTimeout:   5 * time.Second,
				MaxHeaderBytes: 4096,
			},
			wantReadTimeout:  2 * time.Second,
			wantWriteTimeout: 5 * time.Second,
			wantHeaderBytes:  4096,
		},
		{
			name:             "defaults from constructor",
			serverConfig:     config.NewServerConfig("127.0.0.1", 8080),
			wantReadTimeout:  10 * time.Second,
			wantWriteTimeout: 10 * time.Second,
			wantHeaderBytes:  1 << 20,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a := NewApplication().UseConfig(config.Config{Server: tc.serverConfig}).UseServer()

			if a.server.ReadTimeout != tc.wantReadTimeout {
				t.Fatalf("expected read timeout %s, got %s", tc.wantReadTimeout, a.server.ReadTimeout)
			}
			if a.server.WriteTimeout != tc.wantWriteTimeout {
				t.Fatalf("expected write timeout %s, got %s", tc.wantWriteTimeout, a.server.WriteTimeout)
			}
			if a.server.MaxHeaderBytes != tc.wantHeaderBytes {
				t.Fatalf("expected max header bytes %d, got %d", tc.wantHeaderBytes, a.server.MaxHeaderBytes)
			}
		})
	}
}

func TestRouteRegistration_SucceedsAfterFullBootstrap(t *testing.T) {
	a := NewApplication().
		UseConfig(testConfig()).
		UseServer().
		UseServerSecurity(&fakeValidator{})

	if err := a.RegisterGET("/public-get", func(c *gin.Context) { c.String(http.StatusOK, "get") }); err != nil {
		t.Fatalf("register GET: %v", err)
	}

	if err := a.RegisterPOST("/public-post", func(c *gin.Context) { c.String(http.StatusCreated, "post") }); err != nil {
		t.Fatalf("register POST: %v", err)
	}

	if err := a.RegisterProtectedGET("/protected", func(c *gin.Context) { c.String(http.StatusOK, "protected") }); err != nil {
		t.Fatalf("register protected GET: %v", err)
	}

	wGet := httptest.NewRecorder()
	a.handler.ServeHTTP(wGet, httptest.NewRequest(http.MethodGet, "/public-get", nil))
	if wGet.Code != http.StatusOK {
		t.Fatalf("GET /public-get expected 200 got %d", wGet.Code)
	}

	wPost := httptest.NewRecorder()
	a.handler.ServeHTTP(wPost, httptest.NewRequest(http.MethodPost, "/public-post", nil))
	if wPost.Code != http.StatusCreated {
		t.Fatalf("POST /public-post expected 201 got %d", wPost.Code)
	}

	wProtected := httptest.NewRecorder()
	a.handler.ServeHTTP(wProtected, httptest.NewRequest(http.MethodGet, "/protected", nil))
	if wProtected.Code != http.StatusUnauthorized {
		t.Fatalf("GET /protected expected 401 without token got %d", wProtected.Code)
	}
}

func TestRegisterProtectedGET_Enforcement_MissingAndInvalidToken(t *testing.T) {
	validator := &fakeValidator{
		validateFn: func(_ context.Context, raw string) (claims.Claims, error) {
			if raw == "good-token" {
				return claims.Claims{Subject: "user-1"}, nil
			}
			return claims.Claims{}, errors.New("invalid token")
		},
	}

	a := NewApplication().
		UseConfig(testConfig()).
		UseServer().
		UseServerSecurity(validator)

	if err := a.RegisterProtectedGET("/me", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	}); err != nil {
		t.Fatalf("register protected route: %v", err)
	}

	missing := httptest.NewRecorder()
	a.handler.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/me", nil))
	if missing.Code != http.StatusUnauthorized {
		t.Fatalf("missing token expected 401 got %d", missing.Code)
	}

	invalidReq := httptest.NewRequest(http.MethodGet, "/me", nil)
	invalidReq.Header.Set("Authorization", "Bearer bad-token")
	invalid := httptest.NewRecorder()
	a.handler.ServeHTTP(invalid, invalidReq)
	if invalid.Code != http.StatusUnauthorized {
		t.Fatalf("invalid token expected 401 got %d", invalid.Code)
	}
}

func TestOrderingGuards_ReturnExpectedErrors(t *testing.T) {
	t.Run("register get before server", func(t *testing.T) {
		a := NewApplication()
		err := a.RegisterGET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })
		if !errors.Is(err, ErrServerNotReady) {
			t.Fatalf("expected ErrServerNotReady, got %v", err)
		}
	})

	t.Run("register post before server", func(t *testing.T) {
		a := NewApplication()
		err := a.RegisterPOST("/x", func(c *gin.Context) { c.Status(http.StatusOK) })
		if !errors.Is(err, ErrServerNotReady) {
			t.Fatalf("expected ErrServerNotReady, got %v", err)
		}
	})

	t.Run("protected route before security", func(t *testing.T) {
		a := NewApplication().UseConfig(testConfig()).UseServer()
		err := a.RegisterProtectedGET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })
		if !errors.Is(err, ErrSecurityNotReady) {
			t.Fatalf("expected ErrSecurityNotReady, got %v", err)
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		a := NewApplication().UseConfig(testConfig()).UseServer()
		err := a.RegisterGET("   ", func(c *gin.Context) { c.Status(http.StatusOK) })
		if !errors.Is(err, ErrInvalidPath) {
			t.Fatalf("expected ErrInvalidPath, got %v", err)
		}
	})

	t.Run("nil handler", func(t *testing.T) {
		a := NewApplication().UseConfig(testConfig()).UseServer()
		err := a.RegisterGET("/x", nil)
		if !errors.Is(err, ErrNilHandler) {
			t.Fatalf("expected ErrNilHandler, got %v", err)
		}
	})

	t.Run("run before server", func(t *testing.T) {
		a := NewApplication().UseConfig(testConfig())
		err := a.Run()
		if !errors.Is(err, ErrServerNotReady) {
			t.Fatalf("expected ErrServerNotReady, got %v", err)
		}
	})

	t.Run("run listener nil before server", func(t *testing.T) {
		a := NewApplication().UseConfig(testConfig())
		err := a.RunListener(nil)
		if !errors.Is(err, ErrServerNotReady) {
			t.Fatalf("expected ErrServerNotReady, got %v", err)
		}
	})

	t.Run("run listener nil after server", func(t *testing.T) {
		a := NewApplication().UseConfig(testConfig()).UseServer()
		err := a.RunListener(nil)
		if !errors.Is(err, ErrNilListener) {
			t.Fatalf("expected ErrNilListener, got %v", err)
		}
	})
}

func TestEnableHealthReadinessPresets_OptionsAndOrdering(t *testing.T) {
	t.Parallel()

	t.Run("fails before server bootstrap", func(t *testing.T) {
		a := NewApplication()

		err := a.EnableHealthReadinessPresets(HealthReadinessOptions{})
		if !errors.Is(err, ErrServerNotReady) {
			t.Fatalf("expected ErrServerNotReady, got %v", err)
		}
	})

	tests := []struct {
		name    string
		opts    HealthReadinessOptions
		wantErr error
	}{
		{
			name:    "blank health path rejected",
			opts:    HealthReadinessOptions{HealthPath: "   "},
			wantErr: ErrInvalidPresetOptions,
		},
		{
			name:    "blank ready path rejected",
			opts:    HealthReadinessOptions{ReadyPath: "\t"},
			wantErr: ErrInvalidPresetOptions,
		},
		{
			name:    "same path rejected after defaults resolution",
			opts:    HealthReadinessOptions{HealthPath: "/same", ReadyPath: "/same"},
			wantErr: ErrInvalidPresetOptions,
		},
		{
			name: "defaults are accepted",
			opts: HealthReadinessOptions{},
		},
		{
			name: "custom paths are accepted",
			opts: HealthReadinessOptions{HealthPath: "/live", ReadyPath: "/ready"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a := NewApplication().UseConfig(testConfig()).UseServer()

			err := a.EnableHealthReadinessPresets(tc.opts)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected %v, got %v", tc.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestEnableHealthReadinessPresets_ConflictPreflightAndNoPartialRegistration(t *testing.T) {
	t.Parallel()

	t.Run("health conflict fails and does not install readiness", func(t *testing.T) {
		a := NewApplication().UseConfig(testConfig()).UseServer()

		if err := a.RegisterGET("/healthz", func(c *gin.Context) { c.Status(http.StatusAccepted) }); err != nil {
			t.Fatalf("register conflicting route: %v", err)
		}

		err := a.EnableHealthReadinessPresets(HealthReadinessOptions{ReadyPath: "/readyz-alt"})
		if !errors.Is(err, ErrRouteConflict) {
			t.Fatalf("expected ErrRouteConflict, got %v", err)
		}
		if !strings.Contains(err.Error(), "method=GET") || !strings.Contains(err.Error(), "path=\"/healthz\"") {
			t.Fatalf("expected error context with method/path, got %q", err.Error())
		}

		wExisting := httptest.NewRecorder()
		a.handler.ServeHTTP(wExisting, httptest.NewRequest(http.MethodGet, "/healthz", nil))
		if wExisting.Code != http.StatusAccepted {
			t.Fatalf("existing route should remain unchanged, got %d", wExisting.Code)
		}

		wReady := httptest.NewRecorder()
		a.handler.ServeHTTP(wReady, httptest.NewRequest(http.MethodGet, "/readyz-alt", nil))
		if wReady.Code != http.StatusNotFound {
			t.Fatalf("ready route must not be partially registered, got %d", wReady.Code)
		}
	})

	t.Run("ready conflict fails", func(t *testing.T) {
		a := NewApplication().UseConfig(testConfig()).UseServer()

		if err := a.RegisterGET("/readyz", func(c *gin.Context) { c.Status(http.StatusAccepted) }); err != nil {
			t.Fatalf("register conflicting route: %v", err)
		}

		err := a.EnableHealthReadinessPresets(HealthReadinessOptions{})
		if !errors.Is(err, ErrRouteConflict) {
			t.Fatalf("expected ErrRouteConflict, got %v", err)
		}
	})
}

func TestEnableHealthReadinessPresets_HTTPBehavior_DefaultAndCustomPaths(t *testing.T) {
	t.Parallel()

	t.Run("default paths with readiness pass and fail", func(t *testing.T) {
		pass := func() error { return nil }
		fail := func() error { return errors.New("dependency down") }

		aPass := NewApplication().UseConfig(testConfig()).UseServer()
		if err := aPass.EnableHealthReadinessPresets(HealthReadinessOptions{Checks: []ReadinessCheck{pass, pass}}); err != nil {
			t.Fatalf("enable presets (pass): %v", err)
		}

		wHealth := httptest.NewRecorder()
		aPass.handler.ServeHTTP(wHealth, httptest.NewRequest(http.MethodGet, "/healthz", nil))
		if wHealth.Code != http.StatusOK {
			t.Fatalf("/healthz expected 200 got %d", wHealth.Code)
		}

		wReadyPass := httptest.NewRecorder()
		aPass.handler.ServeHTTP(wReadyPass, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if wReadyPass.Code != http.StatusOK {
			t.Fatalf("/readyz expected 200 when checks pass got %d", wReadyPass.Code)
		}

		aFail := NewApplication().UseConfig(testConfig()).UseServer()
		if err := aFail.EnableHealthReadinessPresets(HealthReadinessOptions{Checks: []ReadinessCheck{pass, fail}}); err != nil {
			t.Fatalf("enable presets (fail): %v", err)
		}

		wReadyFail := httptest.NewRecorder()
		aFail.handler.ServeHTTP(wReadyFail, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if wReadyFail.Code != http.StatusServiceUnavailable {
			t.Fatalf("/readyz expected 503 when a check fails got %d", wReadyFail.Code)
		}
	})

	t.Run("custom paths are honored and defaults not duplicated", func(t *testing.T) {
		a := NewApplication().UseConfig(testConfig()).UseServer()
		err := a.EnableHealthReadinessPresets(HealthReadinessOptions{
			HealthPath: "/livez",
			ReadyPath:  "/readyz-custom",
			Checks: []ReadinessCheck{
				func() error { return nil },
			},
		})
		if err != nil {
			t.Fatalf("enable custom presets: %v", err)
		}

		wCustomHealth := httptest.NewRecorder()
		a.handler.ServeHTTP(wCustomHealth, httptest.NewRequest(http.MethodGet, "/livez", nil))
		if wCustomHealth.Code != http.StatusOK {
			t.Fatalf("custom health expected 200 got %d", wCustomHealth.Code)
		}

		wCustomReady := httptest.NewRecorder()
		a.handler.ServeHTTP(wCustomReady, httptest.NewRequest(http.MethodGet, "/readyz-custom", nil))
		if wCustomReady.Code != http.StatusOK {
			t.Fatalf("custom ready expected 200 got %d", wCustomReady.Code)
		}

		wDefaultHealth := httptest.NewRecorder()
		a.handler.ServeHTTP(wDefaultHealth, httptest.NewRequest(http.MethodGet, "/healthz", nil))
		if wDefaultHealth.Code != http.StatusNotFound {
			t.Fatalf("default health path should not be duplicated, got %d", wDefaultHealth.Code)
		}

		wDefaultReady := httptest.NewRecorder()
		a.handler.ServeHTTP(wDefaultReady, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if wDefaultReady.Code != http.StatusNotFound {
			t.Fatalf("default ready path should not be duplicated, got %d", wDefaultReady.Code)
		}
	})

	t.Run("unmet invariant returns 503", func(t *testing.T) {
		a := NewApplication().UseConfig(testConfig()).UseServer()
		if err := a.EnableHealthReadinessPresets(HealthReadinessOptions{}); err != nil {
			t.Fatalf("enable presets: %v", err)
		}

		a.server.Handler = nil

		wReady := httptest.NewRecorder()
		a.handler.ServeHTTP(wReady, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if wReady.Code != http.StatusServiceUnavailable {
			t.Fatalf("/readyz expected 503 when invariant is unmet got %d", wReady.Code)
		}
	})
}

func TestManualRouteRegistration_UnchangedWithoutPresets(t *testing.T) {
	t.Parallel()

	a := NewApplication().UseConfig(testConfig()).UseServer()
	if err := a.RegisterGET("/health", func(c *gin.Context) { c.Status(http.StatusNoContent) }); err != nil {
		t.Fatalf("register manual route: %v", err)
	}

	wManual := httptest.NewRecorder()
	a.handler.ServeHTTP(wManual, httptest.NewRequest(http.MethodGet, "/health", nil))
	if wManual.Code != http.StatusNoContent {
		t.Fatalf("manual route expected 204 got %d", wManual.Code)
	}

	wPresetHealth := httptest.NewRecorder()
	a.handler.ServeHTTP(wPresetHealth, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if wPresetHealth.Code != http.StatusNotFound {
		t.Fatalf("/healthz should not exist without explicit preset opt-in, got %d", wPresetHealth.Code)
	}

	wPresetReady := httptest.NewRecorder()
	a.handler.ServeHTTP(wPresetReady, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if wPresetReady.Code != http.StatusNotFound {
		t.Fatalf("/readyz should not exist without explicit preset opt-in, got %d", wPresetReady.Code)
	}
}

func TestRunListener_ServesRequestAndStopsCleanly(t *testing.T) {
	a := NewApplication().UseConfig(testConfig()).UseServer()
	if err := a.RegisterGET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") }); err != nil {
		t.Fatalf("register health route: %v", err)
	}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.RunListener(l)
	}()

	url := "http://" + l.Addr().String() + "/health"
	resp, err := waitHTTPGet(url, 2*time.Second)
	if err != nil {
		t.Fatalf("request health endpoint: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", resp.StatusCode, string(body))
	}

	if err := a.server.Close(); err != nil {
		t.Fatalf("server close: %v", err)
	}

	select {
	case serveErr := <-errCh:
		if !errors.Is(serveErr, http.ErrServerClosed) {
			t.Fatalf("expected http.ErrServerClosed, got %v", serveErr)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting RunListener to return")
	}
}

func TestRunListener_ReturnsStartupErrors(t *testing.T) {
	a := NewApplication().UseConfig(testConfig()).UseServer()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	if err := l.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}

	err = a.RunListener(l)
	if err == nil {
		t.Fatalf("expected startup error for closed listener")
	}
}

func TestRun_DelegatesToListenAndServeAndPropagatesErrors(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen occupied: %v", err)
	}
	defer occupied.Close()

	tcpAddr, ok := occupied.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("expected TCP listener")
	}

	host := tcpAddr.IP.String()
	if host == "" {
		host = "127.0.0.1"
	}

	a := NewApplication().UseConfig(config.Config{
		Server: config.ServerConfig{Host: host, Port: tcpAddr.Port},
	}).UseServer()

	err = a.Run()
	if err == nil {
		t.Fatalf("expected bind error from ListenAndServe")
	}

	if !strings.Contains(strings.ToLower(err.Error()), "address already in use") {
		t.Fatalf("expected address-in-use error, got %v", err)
	}

	if wantAddr := net.JoinHostPort(host, strconv.Itoa(tcpAddr.Port)); a.server.Addr != wantAddr {
		t.Fatalf("expected server.Addr=%q got %q", wantAddr, a.server.Addr)
	}
}

func TestUseServerSecurityFromConfig(t *testing.T) {
	t.Parallel()

	t.Run("hs256 success", func(t *testing.T) {
		a := NewApplication().UseConfig(config.Config{
			Server: config.NewServerConfig("127.0.0.1", 8080),
			Security: config.NewSecurityConfig(config.NewAuthConfig(
				config.NewJWTConfig("secret", "common-fwk", 15),
				config.NewCookieConfig("session", "example.com", true, true, "Lax"),
				config.NewLoginConfig("owner@example.com"),
				config.NewOAuth2Config(nil),
			)),
		})

		_, err := a.UseServerSecurityFromConfig()
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if !a.securityReady {
			t.Fatalf("expected security to be marked ready")
		}
		if a.validator == nil {
			t.Fatalf("expected validator to be set")
		}
	})

	t.Run("rs256 success", func(t *testing.T) {
		privatePEM := mustRSAPrivatePEM(t)

		jwtCfg := config.NewJWTConfig("", "common-fwk", 15)
		jwtCfg.Algorithm = config.JWTAlgorithmRS256
		jwtCfg.RS256 = config.NewRS256PrivatePEMConfig("rsa-key", privatePEM)

		a := NewApplication().UseConfig(config.Config{
			Server: config.NewServerConfig("127.0.0.1", 8080),
			Security: config.NewSecurityConfig(config.NewAuthConfig(
				jwtCfg,
				config.NewCookieConfig("session", "example.com", true, true, "Lax"),
				config.NewLoginConfig("owner@example.com"),
				config.NewOAuth2Config(nil),
			)),
		})

		_, err := a.UseServerSecurityFromConfig()
		if err != nil {
			t.Fatalf("expected RS256 success, got %v", err)
		}
		if !a.securityReady || a.validator == nil {
			t.Fatalf("expected security wiring to be complete")
		}
	})

	t.Run("invalid config does not partially wire", func(t *testing.T) {
		a := NewApplication().UseConfig(config.Config{
			Server: config.NewServerConfig("127.0.0.1", 8080),
			Security: config.NewSecurityConfig(config.NewAuthConfig(
				config.JWTConfig{Algorithm: config.JWTAlgorithmRS256, Issuer: "common-fwk", TTLMinutes: 15},
				config.NewCookieConfig("session", "example.com", true, true, "Lax"),
				config.NewLoginConfig("owner@example.com"),
				config.NewOAuth2Config(nil),
			)),
		})

		_, err := a.UseServerSecurityFromConfig()
		if err == nil {
			t.Fatalf("expected config-driven security wiring error")
		}
		if a.securityReady {
			t.Fatalf("expected securityReady to remain false on failure")
		}
		if a.validator != nil {
			t.Fatalf("expected validator to remain nil on failure")
		}
	})
}

func TestAccessors_LifecycleMatrix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		bootstrap     func(*Application)
		wantConfig    config.Config
		wantValidator bool
		wantSecReady  bool
	}{
		{
			name:          "pre-init returns explicit non-ready values",
			bootstrap:     func(_ *Application) {},
			wantConfig:    config.Config{},
			wantValidator: false,
			wantSecReady:  false,
		},
		{
			name: "partial-init after UseConfig exposes only config",
			bootstrap: func(a *Application) {
				a.UseConfig(testConfigWithOAuth2Provider())
			},
			wantConfig:    testConfigWithOAuth2Provider(),
			wantValidator: false,
			wantSecReady:  false,
		},
		{
			name: "post-init with direct validator exposes both config and security",
			bootstrap: func(a *Application) {
				a.UseConfig(testConfigWithOAuth2Provider()).UseServerSecurity(&fakeValidator{})
			},
			wantConfig:    testConfigWithOAuth2Provider(),
			wantValidator: true,
			wantSecReady:  true,
		},
		{
			name: "post-init with config-driven security exposes both config and security",
			bootstrap: func(a *Application) {
				a.UseConfig(testConfigWithOAuth2Provider())
				_, err := a.UseServerSecurityFromConfig()
				if err != nil {
					t.Fatalf("expected config-driven security wiring success, got %v", err)
				}
			},
			wantConfig:    testConfigWithOAuth2Provider(),
			wantValidator: true,
			wantSecReady:  true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a := NewApplication()
			tc.bootstrap(a)

			var gotCfg config.Config
			var gotValidator any
			var gotReady bool

			mustNotPanic(t, "GetConfig", func() {
				gotCfg = a.GetConfig()
			})
			mustNotPanic(t, "GetSecurityValidator", func() {
				gotValidator = a.GetSecurityValidator()
			})
			mustNotPanic(t, "IsSecurityReady", func() {
				gotReady = a.IsSecurityReady()
			})

			if !reflect.DeepEqual(gotCfg, tc.wantConfig) {
				t.Fatalf("unexpected config snapshot\nwant=%#v\ngot =%#v", tc.wantConfig, gotCfg)
			}
			if (gotValidator != nil) != tc.wantValidator {
				t.Fatalf("unexpected validator presence: want=%t got=%t", tc.wantValidator, gotValidator != nil)
			}
			if gotReady != tc.wantSecReady {
				t.Fatalf("unexpected security readiness: want=%t got=%t", tc.wantSecReady, gotReady)
			}
		})
	}
}

func TestGetConfig_DefensiveSnapshotImmutability(t *testing.T) {
	t.Parallel()

	a := NewApplication().UseConfig(testConfigWithOAuth2Provider())

	firstSnapshot := a.GetConfig()
	if firstSnapshot.Security.Auth.OAuth2.Providers == nil {
		t.Fatalf("expected providers map in snapshot")
	}

	provider, ok := firstSnapshot.Security.Auth.OAuth2.Providers["github"]
	if !ok {
		t.Fatalf("expected github provider in snapshot")
	}

	provider.Scopes[0] = "mutated-scope"
	firstSnapshot.Security.Auth.OAuth2.Providers["github"] = provider
	firstSnapshot.Security.Auth.OAuth2.Providers["evil"] = config.NewOAuth2ProviderConfig(
		"x",
		"y",
		"https://example.com/auth",
		"https://example.com/token",
		"https://example.com/callback",
		[]string{"scope"},
	)

	secondSnapshot := a.GetConfig()
	provider2, ok := secondSnapshot.Security.Auth.OAuth2.Providers["github"]
	if !ok {
		t.Fatalf("expected github provider in second snapshot")
	}

	if provider2.Scopes[0] != "read:user" {
		t.Fatalf("expected internal scopes to remain unchanged, got %q", provider2.Scopes[0])
	}
	if _, exists := secondSnapshot.Security.Auth.OAuth2.Providers["evil"]; exists {
		t.Fatalf("expected injected provider key to not leak into internal runtime state")
	}

	thirdSnapshot := a.GetConfig()
	secondProvider := secondSnapshot.Security.Auth.OAuth2.Providers["github"]
	thirdProvider := thirdSnapshot.Security.Auth.OAuth2.Providers["github"]
	if &secondProvider.Scopes[0] == &thirdProvider.Scopes[0] {
		t.Fatalf("expected scope slices to be independently copied per read")
	}
}

func TestAccessors_FailedConfigDrivenSecurityRemainsUnavailable(t *testing.T) {
	t.Parallel()

	a := NewApplication().UseConfig(config.Config{
		Server: config.NewServerConfig("127.0.0.1", 8080),
		Security: config.NewSecurityConfig(config.NewAuthConfig(
			config.JWTConfig{Algorithm: config.JWTAlgorithmRS256, Issuer: "common-fwk", TTLMinutes: 15},
			config.NewCookieConfig("session", "example.com", true, true, "Lax"),
			config.NewLoginConfig("owner@example.com"),
			config.NewOAuth2Config(nil),
		)),
	})

	_, err := a.UseServerSecurityFromConfig()
	if err == nil {
		t.Fatalf("expected config-driven security wiring error")
	}

	if a.GetSecurityValidator() != nil {
		t.Fatalf("expected accessor validator to remain nil after failed wiring")
	}
	if a.IsSecurityReady() {
		t.Fatalf("expected accessor security readiness to remain false after failed wiring")
	}
}

func TestDocumentation_AccessorContractSynchronization(t *testing.T) {
	t.Parallel()

	type docSpec struct {
		name string
		path string
	}

	docs := []docSpec{
		{name: "package docs", path: "doc.go"},
		{name: "readme", path: "../README.md"},
		{name: "docs home", path: "../docs/home.md"},
	}

	sharedSignatures := []string{
		"GetConfig() config.Config",
		"GetSecurityValidator() security.Validator",
		"IsSecurityReady() bool",
	}

	for _, doc := range docs {
		doc := doc
		t.Run(doc.name, func(t *testing.T) {
			raw, err := os.ReadFile(doc.path)
			if err != nil {
				t.Fatalf("read %s (%s): %v", doc.name, doc.path, err)
			}

			text := string(raw)
			lower := strings.ToLower(text)

			for _, signature := range sharedSignatures {
				if !strings.Contains(text, signature) {
					t.Fatalf("%s must include accessor signature %q", doc.name, signature)
				}
			}

			if !(strings.Contains(lower, "pre-init") || strings.Contains(lower, "non-init")) {
				t.Fatalf("%s must describe pre-init/non-init accessor behavior", doc.name)
			}
			if !strings.Contains(lower, "zero-value") || !strings.Contains(lower, "nil") || !strings.Contains(lower, "false") {
				t.Fatalf("%s must document pre-init expectations (zero-value config, nil validator, false readiness)", doc.name)
			}

			if !(strings.Contains(lower, "post-init") || strings.Contains(lower, "after security wiring")) {
				t.Fatalf("%s must describe post-init accessor behavior", doc.name)
			}
			if !strings.Contains(lower, "true") {
				t.Fatalf("%s must document ready/true post-init expectations", doc.name)
			}

			hasImmutabilityWording := strings.Contains(lower, "defensive snapshot") || strings.Contains(lower, "deep-cop")
			if !hasImmutabilityWording {
				t.Fatalf("%s must document defensive snapshot/deep-copy immutability", doc.name)
			}

			hasMutationSafetyWording := strings.Contains(lower, "internal runtime state") || strings.Contains(lower, "app internals")
			if !hasMutationSafetyWording {
				t.Fatalf("%s must document that external mutation does not affect internals", doc.name)
			}
		})
	}
}

func TestDocumentation_HealthReadinessPresetContractSynchronization(t *testing.T) {
	t.Parallel()

	type docSpec struct {
		name string
		path string
	}

	docs := []docSpec{
		{name: "package docs", path: "doc.go"},
		{name: "readme", path: "../README.md"},
		{name: "docs home", path: "../docs/home.md"},
	}

	for _, doc := range docs {
		doc := doc
		t.Run(doc.name, func(t *testing.T) {
			raw, err := os.ReadFile(doc.path)
			if err != nil {
				t.Fatalf("read %s (%s): %v", doc.name, doc.path, err)
			}

			text := string(raw)
			lower := strings.ToLower(text)

			if !strings.Contains(text, "EnableHealthReadinessPresets(opts HealthReadinessOptions) error") {
				t.Fatalf("%s must include explicit preset API signature", doc.name)
			}

			if !strings.Contains(text, "/healthz") || !strings.Contains(text, "/readyz") {
				t.Fatalf("%s must document default preset paths /healthz and /readyz", doc.name)
			}

			hasCustomPathBehavior := strings.Contains(lower, "custom") &&
				(strings.Contains(lower, "not duplicated") || strings.Contains(lower, "no implicit duplication"))
			if !hasCustomPathBehavior {
				t.Fatalf("%s must document custom-path behavior without implicit default duplication", doc.name)
			}

			hasReadinessSemantics := strings.Contains(lower, "readiness") &&
				strings.Contains(lower, "200") &&
				strings.Contains(lower, "503")
			if !hasReadinessSemantics {
				t.Fatalf("%s must document readiness 200/503 contract", doc.name)
			}

			hasNoImplicitRegistration := strings.Contains(lower, "no implicit") || strings.Contains(lower, "never auto-registered")
			if !hasNoImplicitRegistration {
				t.Fatalf("%s must document non-goal: no implicit preset registration", doc.name)
			}

			hasNoProviderProbing := strings.Contains(lower, "no provider-specific") || strings.Contains(lower, "provider-specific probing")
			if !hasNoProviderProbing {
				t.Fatalf("%s must document non-goal: no provider-specific probing", doc.name)
			}
		})
	}
}

func mustNotPanic(t *testing.T, name string, fn func()) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("%s panicked: %v", name, r)
		}
	}()

	fn()
}

func waitHTTPGet(url string, timeout time.Duration) (*http.Response, error) {
	deadline := time.Now().Add(timeout)
	var lastErr error

	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		time.Sleep(20 * time.Millisecond)
	}

	return nil, fmt.Errorf("timeout waiting for %s: %w", url, lastErr)
}

func mustRSAPrivatePEM(t *testing.T) string {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}

	der := x509.MarshalPKCS1PrivateKey(priv)
	blk := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}

	return string(pem.EncodeToMemory(blk))
}
