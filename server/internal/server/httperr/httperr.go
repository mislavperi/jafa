// Package httperr centralizes HTTP error responses for controllers.
// All controllers must respond with a sanitized JSON shape ({"error": "<msg>"})
// and log the underlying error server-side. Avoid leaking Go error internals.
package httperr

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Respond writes a sanitized JSON error to the client and logs the real error.
// message should be user-safe (no DB internals). Pass nil err to skip logging.
func Respond(ctx *gin.Context, status int, message string, err error) {
	if err != nil {
		log.Printf("[%s %s] %s: %v", ctx.Request.Method, ctx.Request.URL.Path, message, err)
	}
	ctx.AbortWithStatusJSON(status, gin.H{"error": message})
}

// BadRequest is shorthand for 400.
func BadRequest(ctx *gin.Context, message string, err error) {
	Respond(ctx, http.StatusBadRequest, message, err)
}

// Internal is shorthand for 500.
func Internal(ctx *gin.Context, err error) {
	Respond(ctx, http.StatusInternalServerError, "internal server error", err)
}

// Unauthorized is shorthand for 401.
func Unauthorized(ctx *gin.Context, message string) {
	Respond(ctx, http.StatusUnauthorized, message, nil)
}
