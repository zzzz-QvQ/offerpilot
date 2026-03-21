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
		tokenString := strings.TrimSpace(c.Query("token"))
		authorization := c.GetHeader("Authorization")

		if authorization != "" {
			parts := strings.SplitN(authorization, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Error(c, http.StatusUnauthorized, 40100, "invalid authorization header")
				c.Abort()
				return
			}
			tokenString = strings.TrimSpace(parts[1])
		}

		if tokenString == "" {
			response.Error(c, http.StatusUnauthorized, 40100, "missing authorization token")
			c.Abort()
			return
		}

		claims, err := jwtutil.ParseToken(secret, tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 40102, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Next()
	}
}