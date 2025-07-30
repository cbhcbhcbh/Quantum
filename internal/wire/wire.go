package wire

import (
	"github.com/cbhcbhcbh/Quantum/pkg/common/server"
	"github.com/cbhcbhcbh/Quantum/pkg/config"
	"github.com/cbhcbhcbh/Quantum/pkg/infra"
	"github.com/google/wire"
)

func InitalizeChatServer(name string) (*server.Server, error) {
	wire.Build(
		config.NewConfig,

		infra.NewredisClient,
		infra.NewRedisCacheImpl,
		wire.Bind(new(infra.RedisCache), new(*infra.RedisCacheImpl)),

		infra.NewKafkaPublisher,
		infra.NewKafkaSubscriber,
	)

	return &server.Server{}, nil
}
