package forwarder

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/pkg/common/server"
)

type Router struct {
	grpcServer server.GrpcServer
}

func NewRouter(grpcServer server.GrpcServer) *Router {
	return &Router{grpcServer}
}

func (r *Router) Run() {
	r.grpcServer.Register()
	r.grpcServer.Run()
}
func (r *Router) GracefulStop(ctx context.Context) error {
	return r.grpcServer.GracefulStop()
}
