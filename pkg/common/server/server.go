package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type HttpServer interface {
	RegisterRoutes()
	Run()
	GracefulStop(ctx context.Context) error
}

type GrpcServer interface {
	Register()
	Run()
	GracefulStop() error
}

type Router interface {
	Run()
	GracefulStop(ctx context.Context) error
}

type InfraCloser interface {
	Close() error
}

type Server struct {
	name        string
	router      Router
	infraCloser InfraCloser
}

func NewServer(name string, router Router, infraCloser InfraCloser) *Server {
	return &Server{name, router, infraCloser}
}

func (s *Server) Serve() {
	s.router.Run()

	done := make(chan bool, 1)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.GracefulStop(ctx, done)
	}()

	<-done
}

func (s *Server) GracefulStop(ctx context.Context, done chan bool) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	err := s.router.GracefulStop(ctx)
	if err != nil {
		logger.Error("router graceful stop error", zap.Error(err))
	}

	if err = s.infraCloser.Close(); err != nil {
		logger.Error("infra closer error", zap.Error(err))
	}

	logger.Info("gracefully shutdowned")
	done <- true
}
