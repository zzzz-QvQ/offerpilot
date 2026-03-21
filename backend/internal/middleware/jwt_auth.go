package middleware

import (
	"net/http"
	"offerpilot/backend/internal/pkg/jwtutil"
	"offerpilot/backend/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "currentUserID"

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			response.Error(c, http.StatusUnauthorized, 40100, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, 40100, "invalid authorization header")
			c.Abort()
			return
		}

		claims, err := jwtutil.ParseToken(secret, parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 40102, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Next()
	}
}
