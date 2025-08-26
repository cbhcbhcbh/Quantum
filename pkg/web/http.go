package web

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/cbhcbhcbh/Quantum/pkg/common/log"
	"github.com/cbhcbhcbh/Quantum/pkg/common/middleware"
	"github.com/cbhcbhcbh/Quantum/pkg/config"
)

type HttpServer struct {
	name       string
	logger     log.HttpLog
	svr        *gin.Engine
	httpPort   string
	httpServer *http.Server
}

func NewGinServer(name string, logger log.HttpLog) *gin.Engine {
	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(middleware.LoggingMiddleware(logger))
	return svr
}

func NewHttpServer(name string, logger log.HttpLog, config *config.Config, svr *gin.Engine) *HttpServer {
	return &HttpServer{
		name:     name,
		logger:   logger,
		svr:      svr,
		httpPort: config.Web.Http.Server.Port,
	}
}

func (r *HttpServer) RegisterRoutes() {
	// TODO: Implement static files and javascript files
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

func (r *HttpServer) GracefulStop(ctx context.Context) error {
	return r.httpServer.Shutdown(ctx)
}
