package server

import "context"

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

func (s *Server) Run() {

}

func (s *Server) GracefulStop(ctx context.Context) error {
	if s.infraCloser != nil {
		return s.infraCloser.Close()
	}
	return nil
}
