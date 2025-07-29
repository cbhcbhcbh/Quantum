package wire

import (
	"github.com/cbhcbhcbh/Quantum/pkg/common/server"
	"github.com/cbhcbhcbh/Quantum/pkg/infra"
	"github.com/google/wire"
)

func InitalizeChatServer(name string) (*server.Server, error) {
	wire.Build(
		infra.NewredisClient, 
		infra.NewRedisCacheImpl, 
		wire.Bind(new(infra.RedisCache), new(*infra.RedisCacheImpl)),
	)

	return &server.Server{}, nil
}
