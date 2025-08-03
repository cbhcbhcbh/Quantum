package chat

import (
	"context"

	chatpb "github.com/cbhcbhcbh/Quantum/pkg/proto/chat"
)

func (srv *GrpcServer) CreateChannel(ctx context.Context, req *chatpb.CreateChannelRequest) (*chatpb.CreateChannelResponse, error) {
	// TODO: implement your logic
	return &chatpb.CreateChannelResponse{}, nil
}

func (srv *GrpcServer) AddUserToChannel(ctx context.Context, req *chatpb.AddUserRequest) (*chatpb.AddUserResponse, error) {
	// TODO: implement your logic
	return &chatpb.AddUserResponse{}, nil
}
