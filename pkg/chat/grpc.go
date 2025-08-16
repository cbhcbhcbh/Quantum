package chat

import (
	"net"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/cbhcbhcbh/Quantum/pkg/common/log"
	"github.com/cbhcbhcbh/Quantum/pkg/config"
	chatpb "github.com/cbhcbhcbh/Quantum/pkg/proto/chat"
	"github.com/cbhcbhcbh/Quantum/pkg/transport"
)

type GrpcServer struct {
	grpcPort string
	logger   log.GrpcLog
	s        *grpc.Server
	userSvc  UserService
	chanSvc  ChannelService

	chatpb.UnimplementedChannelServiceServer
	chatpb.UnimplementedUserServiceServer
}

func NewGrpcServer(name string, logger log.GrpcLog, config *config.Config, userSvc UserService, chanSvc ChannelService) *GrpcServer {
	srv := &GrpcServer{
		grpcPort: config.Chat.Grpc.Server.Port,
		logger:   logger,
		userSvc:  userSvc,
		chanSvc:  chanSvc,
	}
	srv.s = transport.InitializeGrpcServer(name, srv.logger)
	return srv
}

func (srv *GrpcServer) Register() {
	chatpb.RegisterChannelServiceServer(srv.s, srv)
	chatpb.RegisterUserServiceServer(srv.s, srv)
}

func (srv *GrpcServer) Run() {
	go func() {
		addr := "0.0.0.0:" + srv.grpcPort
		srv.logger.Info("Starting gRPC server", zap.String("address", addr))
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			srv.logger.Error(err.Error())
			os.Exit(1)
		}
		if err := srv.s.Serve(lis); err != nil {
			srv.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
}

func (srv *GrpcServer) GracefulStop() error {
	srv.s.GetServiceInfo()
	return nil
}

var UserConn *UserClientConn

type UserClientConn struct {
	Conn *grpc.ClientConn
}

func NewUserClientConn(config *config.Config) (*UserClientConn, error) {
	conn, err := transport.InitializeGrpcClient(config.Chat.Grpc.Client.User.Endpoint)
	if err != nil {
		return nil, err
	}
	UserConn = &UserClientConn{
		Conn: conn,
	}
	return UserConn, nil
}

var ForwarderConn *ForwarderClientConn

type ForwarderClientConn struct {
	Conn *grpc.ClientConn
}

func NewForwarderClientConn(config *config.Config) (*ForwarderClientConn, error) {
	conn, err := transport.InitializeGrpcClient(config.Chat.Grpc.Client.Forwarder.Endpoint)
	if err != nil {
		return nil, err
	}
	ForwarderConn = &ForwarderClientConn{
		Conn: conn,
	}
	return ForwarderConn, nil
}
