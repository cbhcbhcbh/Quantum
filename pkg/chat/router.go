package chat

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/pkg/common/server"
)

type Router struct {
	httpServer server.HttpServer
	grpcServer server.GrpcServer
}

func NewRouter(httpServer server.HttpServer, grpcServer server.GrpcServer) *Router {
	return &Router{httpServer, grpcServer}
}

func (r *Router) Run() {
	r.httpServer.RegisterRoutes()
	r.httpServer.Run()

	r.grpcServer.Register()
	r.grpcServer.Run()
}

func (r *Router) GracefulStop(ctx context.Context) error {
	if err := r.grpcServer.GracefulStop(); err != nil {
		return err
	}
	return r.httpServer.GracefulStop(ctx)
}
