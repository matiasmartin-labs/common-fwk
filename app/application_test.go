package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
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
