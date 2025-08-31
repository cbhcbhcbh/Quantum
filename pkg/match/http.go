package match

import (
	"context"
	"net/http"
	"os"

	"github.com/cbhcbhcbh/Quantum/pkg/common/log"
	"github.com/cbhcbhcbh/Quantum/pkg/common/middleware"
	"github.com/cbhcbhcbh/Quantum/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"
)

var (
	MelodyMatch MelodyMatchConn
)

type MelodyMatchConn struct {
	*melody.Melody
}

type HttpServer struct {
	name       string
	logger     log.HttpLog
	svr        *gin.Engine
	mm         MelodyMatchConn
	httpPort   string
	httpServer *http.Server
}

func NewMelodyMatchConn() MelodyMatchConn {
	MelodyMatch = MelodyMatchConn{
		melody.New(),
	}
	return MelodyMatch
}

func NewGinServer(name string, logger log.HttpLog, config *config.Config) *gin.Engine {
	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(middleware.CorsMiddleware())
	svr.Use(middleware.LoggingMiddleware(logger))
	svr.Use(middleware.MaxAllowed(config.Match.Http.Server.MaxConn))
	return svr
}

func NewHttpServer(name string, logger log.HttpLog, config *config.Config, svr *gin.Engine, mm MelodyMatchConn) *HttpServer {
	return &HttpServer{
		name:     name,
		logger:   logger,
		svr:      svr,
		mm:       mm,
		httpPort: config.Match.Http.Server.Port,
	}
}

func (h *HttpServer) RegisterRoutes() {
	// TODO: Add prometheus metrics endpoint
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
	err := MelodyMatch.Close()
	if err != nil {
		return err
	}
	err = h.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
