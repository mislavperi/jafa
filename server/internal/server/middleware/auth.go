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
		userID := session.Get(SessionUserIDKey)
		if userID == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}
		var id int64
		switch v := userID.(type) {
		case int64:
			id = v
		case int:
			id = int64(v)
		case int32:
			id = int64(v)
		case float64:
			id = int64(v)
		default:
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
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
