package gin

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const bearerPrefix = "Bearer "

// extractToken resolves the raw token from the request.
// It checks the Authorization header first (Bearer scheme), then falls back
// to the named cookie. Returns an empty string when neither source provides a token.
func extractToken(c *gin.Context, headerName, cookieName string) string {
	if header := c.GetHeader(headerName); header != "" {
		if strings.HasPrefix(header, bearerPrefix) {
			return strings.TrimPrefix(header, bearerPrefix)
		}
		// Non-Bearer scheme — treat as missing (not invalid).
		return ""
	}

	if cookie, err := c.Cookie(cookieName); err == nil && cookie != "" {
		return cookie
	}

	return ""
}
