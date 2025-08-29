package forwarder

import (
	"context"

	forwarderpb "github.com/cbhcbhcbh/Quantum/pkg/proto/forwarder"
)

func (srv *GrpcServer) RegisterChannelSession(ctx context.Context, req *forwarderpb.RegisterChannelSessionRequest) (*forwarderpb.RegisterChannelSessionResponse, error) {
	// TODO: register session to forward service
	return &forwarderpb.RegisterChannelSessionResponse{}, nil
}

func (srv *GrpcServer) RemoveChannelSession(ctx context.Context, req *forwarderpb.RemoveChannelSessionRequest) (*forwarderpb.RemoveChannelSessionResponse, error) {
	// TODO: remove session to forward service
	return &forwarderpb.RemoveChannelSessionResponse{}, nil
}
