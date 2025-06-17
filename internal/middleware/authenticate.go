package middleware

import (
	"net/http"
	"strings"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/cbhcbhcbh/Quantum/pkg/jwt"
	"github.com/cbhcbhcbh/Quantum/pkg/response"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.ErrorResponse(http.StatusUnauthorized, "No authHeader").SetHttpCode(http.StatusUnauthorized).ToJson(c)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorResponse(http.StatusUnauthorized, "No Bearer").SetHttpCode(http.StatusUnauthorized).ToJson(c)
			c.Abort()
			return
		}

		claims, err := jwt.NewJWT().ParseToken(parts[1])
		if err != nil {
			response.ErrorResponse(http.StatusUnauthorized, err.Error()).SetHttpCode(http.StatusUnauthorized).WriteTo(c)
			c.Abort()
			return
		}

		c.Set(known.XIdKey, claims.ID)
		c.Set(known.XUidKey, claims.UID)
		c.Set(known.XUsernameKey, claims.Name)
		c.Set(known.XEmailKey, claims.Email)

		c.Next()
	}
}
