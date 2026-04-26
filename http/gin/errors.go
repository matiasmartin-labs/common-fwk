package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is the JSON body returned on authentication failures.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeError aborts the request and writes a JSON error response with HTTP 401.
func writeError(c *gin.Context, code, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
		Code:    code,
		Message: msg,
	})
}
