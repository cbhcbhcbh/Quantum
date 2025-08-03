package transport

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/cbhcbhcbh/Quantum/pkg/common/log"
)

var (
	ServiceIdHeader string = "Service-Id"
)

func interceptorLogger(l log.GrpcLog) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg, zap.Any("fields", fields))
		case logging.LevelInfo:
			l.Info(msg, zap.Any("fields", fields))
		case logging.LevelWarn:
			l.Warn(msg, zap.Any("fields", fields))
		case logging.LevelError:
			l.Error(msg, zap.Any("fields", fields))
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func InitializeGrpcServer(name string, logger log.GrpcLog) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 8),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second,
			MaxConnectionAge:      600 * time.Second,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  5 * time.Second,
			Timeout:               1 * time.Second,
		}),
	}

	srv := grpc.NewServer(opts...)
	return srv
}

func InitializeGrpcClient(svcHost string) (*grpc.ClientConn, error) {
	scheme := "dns"

	client, err := grpc.NewClient(
		fmt.Sprintf("%s:///%s", scheme, svcHost),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewGrpcEndpoint(conn *grpc.ClientConn, serviceName, method string, grpcReply interface{}) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		serviceName,
		method,
		func(_ context.Context, req interface{}) (interface{}, error) { return req, nil },
		func(_ context.Context, resp interface{}) (interface{}, error) { return resp, nil },
		grpcReply,
	).Endpoint()
}
