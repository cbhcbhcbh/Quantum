//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package wire

import (
	"github.com/cbhcbhcbh/Quantum/pkg/chat"
	"github.com/cbhcbhcbh/Quantum/pkg/common/log"
	"github.com/cbhcbhcbh/Quantum/pkg/common/server"
	"github.com/cbhcbhcbh/Quantum/pkg/common/sonyflake"
	"github.com/cbhcbhcbh/Quantum/pkg/config"
	"github.com/cbhcbhcbh/Quantum/pkg/infra"
	"github.com/cbhcbhcbh/Quantum/pkg/user"
	"github.com/cbhcbhcbh/Quantum/pkg/web"
	"github.com/google/wire"
)

func InitializeWebServer(name string) (*server.Server, error) {
	wire.Build(
		config.NewConfig,
		log.NewHttpLog,

		web.NewGinServer,

		web.NewHttpServer,
		wire.Bind(new(server.HttpServer), new(*web.HttpServer)),
		web.NewRouter,
		wire.Bind(new(server.Router), new(*web.Router)),
		web.NewInfraCloser,
		wire.Bind(new(server.InfraCloser), new(*web.InfraCloser)),
		server.NewServer,
	)
	return &server.Server{}, nil
}

func InitializeChatServer(name string) (*server.Server, error) {
	wire.Build(
		config.NewConfig,
		log.NewHttpLog,
		log.NewGrpcLog,

		infra.NewRedisClient,
		infra.NewRedisCacheImpl,
		wire.Bind(new(infra.RedisCache), new(*infra.RedisCacheImpl)),

		infra.NewKafkaPublisher,
		infra.NewKafkaSubscriber,
		infra.NewSimpleRouter,

		infra.NewCassandraSession,

		chat.NewUserClientConn,
		chat.NewForwarderClientConn,

		chat.NewUserRepoImpl,
		wire.Bind(new(chat.UserRepo), new(*chat.UserRepoImpl)),
		chat.NewMessageRepoImpl,
		wire.Bind(new(chat.MessageRepo), new(*chat.MessageRepoImpl)),
		chat.NewChannelRepoImpl,
		wire.Bind(new(chat.ChannelRepo), new(*chat.ChannelRepoImpl)),
		chat.NewForwardRepoImpl,
		wire.Bind(new(chat.ForwardRepo), new(*chat.ForwardRepoImpl)),

		chat.NewUserRepoCacheImpl,
		wire.Bind(new(chat.UserRepoCache), new(*chat.UserRepoCacheImpl)),
		chat.NewMessageRepoCacheImpl,
		wire.Bind(new(chat.MessageRepoCache), new(*chat.MessageRepoCacheImpl)),
		chat.NewChannelRepoCacheImpl,
		wire.Bind(new(chat.ChannelRepoCache), new(*chat.ChannelRepoCacheImpl)),

		chat.NewMessageSubscriber,

		sonyflake.NewSonyFlake,

		chat.NewUserServiceImpl,
		wire.Bind(new(chat.UserService), new(*chat.UserServiceImpl)),
		chat.NewMessageServiceImpl,
		wire.Bind(new(chat.MessageService), new(*chat.MessageServiceImpl)),
		chat.NewChannelServiceImpl,
		wire.Bind(new(chat.ChannelService), new(*chat.ChannelServiceImpl)),
		chat.NewForwardServiceImpl,
		wire.Bind(new(chat.ForwardService), new(*chat.ForwardServiceImpl)),

		chat.NewMelodyChatConn,

		chat.NewGinServer,

		chat.NewHttpServer,
		wire.Bind(new(server.HttpServer), new(*chat.HttpServer)),
		chat.NewGrpcServer,
		wire.Bind(new(server.GrpcServer), new(*chat.GrpcServer)),
		chat.NewRouter,
		wire.Bind(new(server.Router), new(*chat.Router)),
		chat.NewInfraCloser,
		wire.Bind(new(server.InfraCloser), new(*chat.InfraCloser)),
		server.NewServer,
	)

	return &server.Server{}, nil
}

func InitializeUserServer(name string) (*server.Server, error) {
	wire.Build(
		config.NewConfig,
		log.NewHttpLog,
		log.NewGrpcLog,

		infra.NewRedisClient,
		infra.NewRedisCacheImpl,
		wire.Bind(new(infra.RedisCache), new(*infra.RedisCacheImpl)),

		user.NewUserRepoImpl,
		wire.Bind(new(user.UserRepo), new(*user.UserRepoImpl)),

		sonyflake.NewSonyFlake,

		user.NewUserServiceImpl,
		wire.Bind(new(user.UserService), new(*user.UserServiceImpl)),

		user.NewGinServer,

		user.NewHttpServer,
		wire.Bind(new(server.HttpServer), new(*user.HttpServer)),
		user.NewGrpcServer,
		wire.Bind(new(server.GrpcServer), new(*user.GrpcServer)),
		user.NewRouter,
		wire.Bind(new(server.Router), new(*user.Router)),
		user.NewInfraCloser,
		wire.Bind(new(server.InfraCloser), new(*user.InfraCloser)),
		server.NewServer,
	)
	return &server.Server{}, nil
}
