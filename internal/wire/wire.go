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
	"github.com/google/wire"
)

func InitializeChatServer(name string) (*server.Server, error) {
	wire.Build(
		config.NewConfig,
		log.NewHttpLog,
		log.NewGrpcLog,

		infra.NewredisClient,
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
