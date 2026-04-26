package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/matiasmartin-labs/common-fwk/security/claims"
)

// SetClaims stores validated claims in the gin.Context under key.
func SetClaims(c *gin.Context, key string, cl claims.Claims) {
	c.Set(key, cl)
}

// GetClaims retrieves claims from the gin.Context stored under key.
// Returns (claims.Claims{}, false) when absent or if the stored value is
// not of type claims.Claims.
func GetClaims(c *gin.Context, key string) (claims.Claims, bool) {
	v, exists := c.Get(key)
	if !exists {
		return claims.Claims{}, false
	}

	cl, ok := v.(claims.Claims)
	if !ok {
		return claims.Claims{}, false
	}

	return cl, true
}
