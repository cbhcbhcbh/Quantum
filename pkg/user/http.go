package user

import (
	"context"
	"net/http"
	"os"

	"github.com/cbhcbhcbh/Quantum/pkg/common/cookie"
	"github.com/cbhcbhcbh/Quantum/pkg/common/known"
	"github.com/cbhcbhcbh/Quantum/pkg/common/log"
	"github.com/cbhcbhcbh/Quantum/pkg/common/middleware"
	"github.com/cbhcbhcbh/Quantum/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type HttpServer struct {
	name       string
	logger     log.HttpLog
	svr        *gin.Engine
	httpPort   string
	httpServer *http.Server
	userSvc    UserService

	googleOauthConfig *oauth2.Config
	oauthCookieConfig config.CookieConfig
	authCookieConfig  config.CookieConfig
}

func NewGinServer(name string, logger log.HttpLog, config *config.Config) *gin.Engine {
	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(middleware.CorsMiddleware())
	svr.Use(middleware.LoggingMiddleware(logger))

	return svr
}

func NewHttpServer(name string, logger log.HttpLog, config *config.Config, svr *gin.Engine) *HttpServer {
	return &HttpServer{
		name:     name,
		logger:   logger,
		svr:      svr,
		httpPort: config.User.Http.Server.Port,

		googleOauthConfig: &oauth2.Config{
			ClientID:     config.User.OAuth.Google.ClientID,
			ClientSecret: config.User.OAuth.Google.ClientSecret,
			RedirectURL:  config.User.OAuth.Google.RedirectUrl,
			Scopes:       config.User.OAuth.Google.Scopes,
			Endpoint:     google.Endpoint,
		},
		oauthCookieConfig: config.User.OAuth.Cookie,
		authCookieConfig:  config.User.Auth.Cookie,
	}
}

func (h *HttpServer) CookieAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := cookie.GetCookie(c, known.SessionIdCookieName)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, err := h.userSvc.GetUserIDBySession(c.Request.Context(), sid)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), known.UserKey, userID))
		c.Next()
	}
}

func (h *HttpServer) RegisterRoutes() {
	userGroup := h.svr.Group("/api/user")
	{
		userGroup.POST("", h.CreateLocalUser)

		cookieAuthGroup := userGroup.Group("")
		cookieAuthGroup.Use(h.CookieAuth())
		cookieAuthGroup.GET("", h.GetUser)
		cookieAuthGroup.GET("/me", h.GetUserMe)

		userGroup.GET("/oauth2/google/login", h.OAuthGoogleLogin)
		userGroup.GET("/oauth2/google/callback", h.OAuthGoogleCallback)
	}
}

func (h *HttpServer) Run() {
	go func() {
		addr := ":" + h.httpPort
		h.httpServer = &http.Server{
			Addr:    addr,
			Handler: h.svr,
		}
		h.logger.Info("http server listening", zap.String("addr", addr))
		err := h.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			h.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
}

func (h *HttpServer) GracefulStop(ctx context.Context) error {
	return h.httpServer.Shutdown(ctx)
}
