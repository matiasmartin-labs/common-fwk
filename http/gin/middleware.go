package gin

import (
	"github.com/gin-gonic/gin"
	fwkerrors "github.com/matiasmartin-labs/common-fwk/errors"
	"github.com/matiasmartin-labs/common-fwk/security"
)

const (
	defaultHeaderName  = "Authorization"
	defaultCookieName  = "token"
	defaultContextKey  = "claims"
	msgTokenMissing    = "authentication token is missing"
	msgTokenInvalid    = "authentication token is invalid"
)

type options struct {
	authEnabled bool
	headerName  string
	cookieName  string
	contextKey  string
}

// Option configures the auth middleware.
type Option func(*options)

// WithAuthEnabled controls whether authentication is enforced.
// Pass false to allow all requests through without token checks.
func WithAuthEnabled(enabled bool) Option {
	return func(o *options) {
		o.authEnabled = enabled
	}
}

// WithHeaderName sets the HTTP header name used to carry the Bearer token.
func WithHeaderName(name string) Option {
	return func(o *options) {
		o.headerName = name
	}
}

// WithCookieName sets the cookie name used as the fallback token source.
func WithCookieName(name string) Option {
	return func(o *options) {
		o.cookieName = name
	}
}

// WithContextKey sets the gin.Context key used to store validated claims.
func WithContextKey(key string) Option {
	return func(o *options) {
		o.contextKey = key
	}
}

// NewAuthMiddleware returns a Gin handler that authenticates requests using validator.
//
// Flow:
//  1. If auth is disabled (WithAuthEnabled(false)), the request passes through.
//  2. The token is extracted from headerName (Bearer) then cookieName.
//  3. Missing token → 401 with code auth_token_missing.
//  4. Validation error → 401 with code auth_token_invalid.
//  5. Success → claims stored under contextKey and request continues.
func NewAuthMiddleware(validator security.Validator, opts ...Option) gin.HandlerFunc {
	o := &options{
		authEnabled: true,
		headerName:  defaultHeaderName,
		cookieName:  defaultCookieName,
		contextKey:  defaultContextKey,
	}
	for _, opt := range opts {
		opt(o)
	}

	return func(c *gin.Context) {
		if !o.authEnabled {
			c.Next()
			return
		}

		token := extractToken(c, o.headerName, o.cookieName)
		if token == "" {
			writeError(c, fwkerrors.CodeTokenMissing, msgTokenMissing)
			return
		}

		cl, err := validator.Validate(c.Request.Context(), token)
		if err != nil {
			writeError(c, fwkerrors.CodeTokenInvalid, msgTokenInvalid)
			return
		}

		SetClaims(c, o.contextKey, cl)
		c.Next()
	}
}
