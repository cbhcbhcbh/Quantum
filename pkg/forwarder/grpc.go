package forwarder

import (
	"net"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/cbhcbhcbh/Quantum/pkg/common/log"
	"github.com/cbhcbhcbh/Quantum/pkg/config"
	forwarderpb "github.com/cbhcbhcbh/Quantum/pkg/proto/forwarder"
	"github.com/cbhcbhcbh/Quantum/pkg/transport"
)

type GrpcServer struct {
	grpcPort      string
	logger        log.GrpcLog
	s             *grpc.Server
	forwardSvc    ForwardService
	msgSubscriber *MessageSubscriber

	forwarderpb.UnimplementedForwardServiceServer
}

func NewGrpcServer(name string, logger log.GrpcLog, config *config.Config, forwardSvc ForwardService, msgSubscriber *MessageSubscriber) *GrpcServer {
	srv := &GrpcServer{
		grpcPort:      config.Forwarder.Grpc.Server.Port,
		logger:        logger,
		forwardSvc:    forwardSvc,
		msgSubscriber: msgSubscriber,
	}
	srv.s = transport.InitializeGrpcServer(name, srv.logger)
	return srv
}

func (srv *GrpcServer) Register() {
	forwarderpb.RegisterForwardServiceServer(srv.s, srv)
}

func (srv *GrpcServer) Run() {
	go func() {
		addr := "0.0.0.0:" + srv.grpcPort
		srv.logger.Info("grpc server listening", zap.String("address", addr))
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
	go func() {
		err := srv.msgSubscriber.Run()
		if err != nil {
			srv.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
}

func (srv *GrpcServer) GracefulStop() error {
	srv.s.GracefulStop()
	return srv.msgSubscriber.GracefulStop()
}
