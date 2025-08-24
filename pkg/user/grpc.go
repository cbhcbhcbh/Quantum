package user

import (
	"net"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/cbhcbhcbh/Quantum/pkg/common/log"
	"github.com/cbhcbhcbh/Quantum/pkg/config"
	userpb "github.com/cbhcbhcbh/Quantum/pkg/proto/user"
	"github.com/cbhcbhcbh/Quantum/pkg/transport"
)

type GrpcServer struct {
	grpcPort string
	logger   log.GrpcLog
	s        *grpc.Server
	userSvc  UserService

	userpb.UnimplementedUserServiceServer
}

func NewGrpcServer(name string, logger log.GrpcLog, config *config.Config, userSvc UserService) *GrpcServer {
	srv := &GrpcServer{
		grpcPort: config.User.Grpc.Server.Port,
		logger:   logger,
		userSvc:  userSvc,
	}
	srv.s = transport.InitializeGrpcServer(name, srv.logger)
	return srv
}

func (srv *GrpcServer) Register() {
	userpb.RegisterUserServiceServer(srv.s, srv)
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
