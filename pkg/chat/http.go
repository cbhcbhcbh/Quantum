package chat

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
	sessCidKey = "sesscid"

	MelodyChat MelodyChatConn
)

type MelodyChatConn struct {
	*melody.Melody
}

type HttpServer struct {
	name          string
	logger        log.HttpLog
	svr           *gin.Engine
	mc            MelodyChatConn
	httpPort      string
	httpServer    *http.Server
	msgSubscriber *MessageSubscriber
	userSvc       UserService
}

func NewMelodyChatConn(config *config.Config) MelodyChatConn {
	m := melody.New()
	m.Config.MaxMessageSize = config.Chat.Message.MaxSizeByte
	MelodyChat = MelodyChatConn{m}
	return MelodyChat
}

func NewGinServer(name string, logger log.HttpLog, config *config.Config) *gin.Engine {
	svr := gin.New()
	svr.Use(gin.Recovery())
	svr.Use(middleware.CorsMiddleware())
	svr.Use(middleware.LoggingMiddleware(logger))
	svr.Use(middleware.MaxAllowed(config.Chat.Http.Server.MaxConn))

	return svr
}

func NewHttpServer(name string, logger log.HttpLog, config *config.Config, svr *gin.Engine, mc MelodyChatConn, msgSubscriber *MessageSubscriber, userSvc UserService) *HttpServer {
	// TODO: 配置 jwt

	return &HttpServer{
		name:          name,
		logger:        logger,
		svr:           svr,
		mc:            mc,
		httpPort:      config.Chat.Http.Server.Port,
		msgSubscriber: msgSubscriber,
		userSvc:       userSvc,
	}
}

func (h *HttpServer) RegisterRoutes() {
	h.msgSubscriber.RegisterHandler()

	chatGroup := h.svr.Group("/api/chat")
	{
		chatGroup.GET("", h.StartChat)

		usersGroup := chatGroup.Group("/users")
		usersGroup.Use(middleware.Auth())
		{
			usersGroup.GET("", h.GetChannelUsers)
			usersGroup.GET("/online", h.GetOnlineUsers)
		}

		channelGroup := chatGroup.Group("/channel")
		channelGroup.Use(middleware.Auth())
		{
			channelGroup.GET("/messages", h.ListMessages)
			channelGroup.DELETE("", h.DeleteChannel)
		}
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

	go func() {
		if err := h.msgSubscriber.Run(); err != nil {
			h.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
}

func (h *HttpServer) GracefulStop(ctx context.Context) error {
	if err := MelodyChat.Close(); err != nil {
		return err
	}

	if err := h.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	if err := h.msgSubscriber.GracefulStop(); err != nil {
		return err
	}

	return nil
}
