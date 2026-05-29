package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const SessionUserIDKey = "user_id"
const ContextUserIDKey = "user_id"

func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		// The session only ever stores user.Id (int64), so a direct type
		// assertion is enough — no need to handle other numeric types.
		id, ok := session.Get(SessionUserIDKey).(int64)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}
		ctx.Set(ContextUserIDKey, id)
		ctx.Next()
	}
}

func CurrentUserID(ctx *gin.Context) (int64, bool) {
	v, ok := ctx.Get(ContextUserIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}
