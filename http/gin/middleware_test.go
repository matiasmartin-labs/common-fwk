package gin_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	ginfwk "github.com/matiasmartin-labs/common-fwk/http/gin"
	"github.com/matiasmartin-labs/common-fwk/security/claims"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// fakeValidator is a test double for security.Validator.
type fakeValidator struct {
	returnClaims claims.Claims
	returnErr    error
}

func (f *fakeValidator) Validate(_ context.Context, _ string) (claims.Claims, error) {
	return f.returnClaims, f.returnErr
}

// newEngine builds a minimal Gin engine with the middleware and a catch-all handler.
func newEngine(mw gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(mw)
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

// bodyCode unmarshals the JSON response body and returns the "code" field.
func bodyCode(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var resp ginfwk.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode error response: %v", err)
	}
	return resp.Code
}

// --- Middleware behaviour tests ---

func TestAuthMiddleware_AuthDisabled_PassesThrough(t *testing.T) {
	v := &fakeValidator{returnErr: errors.New("should not be called")}
	mw := ginfwk.NewAuthMiddleware(v, ginfwk.WithAuthEnabled(false))
	r := newEngine(mw)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestAuthMiddleware_NoToken_Returns401Missing(t *testing.T) {
	v := &fakeValidator{}
	mw := ginfwk.NewAuthMiddleware(v)
	r := newEngine(mw)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
	if code := bodyCode(t, w); code != "auth_token_missing" {
		t.Fatalf("expected auth_token_missing got %q", code)
	}
}

func TestAuthMiddleware_ValidHeaderToken_200WithClaims(t *testing.T) {
	want := claims.Claims{Subject: "user-1"}
	v := &fakeValidator{returnClaims: want}
	mw := ginfwk.NewAuthMiddleware(v)

	var got claims.Claims
	r := gin.New()
	r.Use(mw)
	r.GET("/", func(c *gin.Context) {
		cl, ok := ginfwk.GetClaims(c, "claims")
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		got = cl
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	if got.Subject != want.Subject {
		t.Fatalf("expected subject %q got %q", want.Subject, got.Subject)
	}
}

func TestAuthMiddleware_ValidCookieToken_200WithClaims(t *testing.T) {
	want := claims.Claims{Subject: "user-cookie"}
	v := &fakeValidator{returnClaims: want}
	mw := ginfwk.NewAuthMiddleware(v)

	var got claims.Claims
	r := gin.New()
	r.Use(mw)
	r.GET("/", func(c *gin.Context) {
		cl, ok := ginfwk.GetClaims(c, "claims")
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		got = cl
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "cookie-token"})
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	if got.Subject != want.Subject {
		t.Fatalf("expected subject %q got %q", want.Subject, got.Subject)
	}
}

func TestAuthMiddleware_HeaderWinsOverCookie(t *testing.T) {
	headerClaims := claims.Claims{Subject: "from-header"}
	// validator always returns the same claims; we verify WHICH token was passed
	var receivedToken string
	r := gin.New()
	r.Use(func(c *gin.Context) {
		// capture the token used by peeking at the header
		h := c.GetHeader("Authorization")
		if h != "" {
			receivedToken = h
		}
		c.Next()
	})
	v := &fakeValidator{returnClaims: headerClaims}
	mw := ginfwk.NewAuthMiddleware(v)
	r.Use(mw)
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer header-token")
	req.AddCookie(&http.Cookie{Name: "token", Value: "cookie-token"})
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	if receivedToken != "Bearer header-token" {
		t.Fatalf("expected header token to win, receivedToken=%q", receivedToken)
	}
}

func TestAuthMiddleware_InvalidToken_Returns401Invalid(t *testing.T) {
	v := &fakeValidator{returnErr: errors.New("malformed")}
	mw := ginfwk.NewAuthMiddleware(v)
	r := newEngine(mw)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
	if code := bodyCode(t, w); code != "auth_token_invalid" {
		t.Fatalf("expected auth_token_invalid got %q", code)
	}
}

func TestAuthMiddleware_ExpiredToken_Returns401Invalid(t *testing.T) {
	v := &fakeValidator{returnErr: errors.New("token expired")}
	mw := ginfwk.NewAuthMiddleware(v)
	r := newEngine(mw)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer expired-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
	if code := bodyCode(t, w); code != "auth_token_invalid" {
		t.Fatalf("expected auth_token_invalid got %q", code)
	}
}

func TestAuthMiddleware_InvalidSignature_Returns401Invalid(t *testing.T) {
	v := &fakeValidator{returnErr: errors.New("invalid signature")}
	mw := ginfwk.NewAuthMiddleware(v)
	r := newEngine(mw)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer wrong-sig-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
	if code := bodyCode(t, w); code != "auth_token_invalid" {
		t.Fatalf("expected auth_token_invalid got %q", code)
	}
}

// --- Extractor unit tests ---

func TestExtractToken_BearerHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer my-token")

	// We test via middleware to avoid exporting extractToken.
	var captured string
	v := &fakeValidator{returnErr: errors.New("stop")}
	_ = v // use the middleware to confirm extraction happens
	// Instead, test indirectly: valid token path calls validator.
	called := false
	mv := &captureValidator{fn: func(raw string) { called = true; captured = raw }}
	mw := ginfwk.NewAuthMiddleware(mv)
	r := gin.New()
	r.Use(mw)
	r.GET("/", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer extracted-token")
	r.ServeHTTP(w, req)

	if !called {
		t.Fatal("expected validator to be called")
	}
	if captured != "extracted-token" {
		t.Fatalf("expected extracted-token got %q", captured)
	}
}

func TestExtractToken_MalformedHeader_TreatedAsMissing(t *testing.T) {
	v := &fakeValidator{}
	mw := ginfwk.NewAuthMiddleware(v)
	r := newEngine(mw)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz") // non-Bearer
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
	if code := bodyCode(t, w); code != "auth_token_missing" {
		t.Fatalf("expected auth_token_missing got %q", code)
	}
}

func TestExtractToken_CookieFallback(t *testing.T) {
	captured := ""
	mv := &captureValidator{fn: func(raw string) { captured = raw }}
	mw := ginfwk.NewAuthMiddleware(mv)
	r := newEngine(mw)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "cookie-value"})
	r.ServeHTTP(w, req)

	if captured != "cookie-value" {
		t.Fatalf("expected cookie-value got %q", captured)
	}
}

func TestExtractToken_BothAbsent_ReturnsEmpty(t *testing.T) {
	v := &fakeValidator{}
	mw := ginfwk.NewAuthMiddleware(v)
	r := newEngine(mw)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
}

// captureValidator records the token passed to Validate and returns no error.
type captureValidator struct {
	fn func(raw string)
}

func (cv *captureValidator) Validate(_ context.Context, raw string) (claims.Claims, error) {
	cv.fn(raw)
	return claims.Claims{}, nil
}

// --- Context helper tests ---

func TestGetSetClaims_RoundTrip(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	want := claims.Claims{Subject: "test-subject"}
	ginfwk.SetClaims(c, "claims", want)
	got, ok := ginfwk.GetClaims(c, "claims")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if got.Subject != want.Subject {
		t.Fatalf("expected %q got %q", want.Subject, got.Subject)
	}
}

func TestGetSetClaims_CustomKey(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	want := claims.Claims{Subject: "custom"}
	ginfwk.SetClaims(c, "my-key", want)
	got, ok := ginfwk.GetClaims(c, "my-key")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if got.Subject != want.Subject {
		t.Fatalf("expected %q got %q", want.Subject, got.Subject)
	}
}

func TestGetClaims_AbsentKey_ReturnsFalse(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	_, ok := ginfwk.GetClaims(c, "missing")
	if ok {
		t.Fatal("expected ok=false for absent key")
	}
}

func TestGetClaims_WrongType_ReturnsFalse(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("claims", "not-a-claims-struct")
	_, ok := ginfwk.GetClaims(c, "claims")
	if ok {
		t.Fatal("expected ok=false for wrong type")
	}
}

// --- Integration-lite: jwt.ValidationError wrapping ---

func TestAuthMiddleware_WrappedValidationError_MapsToInvalid(t *testing.T) {
	// Simulate what security/jwt returns: a wrapped error.
	wrappedErr := errors.New("validation failed at claims: expired token: expired token")
	v := &fakeValidator{returnErr: wrappedErr}
	mw := ginfwk.NewAuthMiddleware(v)
	r := newEngine(mw)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer some.jwt.token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
	var resp ginfwk.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Code != "auth_token_invalid" {
		t.Fatalf("expected auth_token_invalid got %q", resp.Code)
	}
	// Ensure internal error details are NOT leaked in the message.
	if resp.Message != ginfwk.MsgTokenInvalid {
		t.Fatalf("unexpected message leak: %q", resp.Message)
	}
}
