package config

import (
	"context"

	"github.com/cbhcbhcbh/Quantum/internal/service/queue"
	"github.com/spf13/viper"
)

func StartConsumer() {
	ctx := context.Background()
	addr := viper.GetStringSlice("kafka.addr")

	go func() {
		queue.ConsumerMessage(ctx, addr)
	}()
}
