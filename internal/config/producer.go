package config

import (
	"github.com/cbhcbhcbh/Quantum/pkg/kafka"
	"github.com/spf13/viper"
)

func InitKafkaProducer() {
	addr := viper.GetStringSlice("kafka.addr")
	kafka.P = kafka.NewProducer(addr)
}
