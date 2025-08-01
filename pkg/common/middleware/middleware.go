package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"go.uber.org/zap"

	"github.com/cbhcbhcbh/Quantum/pkg/common/jwt"
	"github.com/cbhcbhcbh/Quantum/pkg/common/known"
	"github.com/cbhcbhcbh/Quantum/pkg/common/log"
	"github.com/cbhcbhcbh/Quantum/pkg/common/response"
	"github.com/cbhcbhcbh/Quantum/pkg/common/util"
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

func CorsMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", known.JWTAuthHeader},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(config)
}

func LoggingMiddleware(logger log.HttpLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := util.GetDurationInMillseconds(start)

		logger.Info("",
			zap.Float64("duration_ms", duration),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.RequestURI),
			zap.Int("status", c.Writer.Status()),
			zap.String("referrer", c.Request.Referer()),
		)
	}
}

func MaxAllowed(n int64) gin.HandlerFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(c *gin.Context) {
		acquire()
		defer release()
		c.Next()
	}
}
