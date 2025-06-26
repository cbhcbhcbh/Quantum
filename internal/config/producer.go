package config

import (
	"strings"

	"github.com/cbhcbhcbh/Quantum/pkg/kafka"
	"github.com/spf13/viper"
)

func InitKafkaProducer() {
	hosts := viper.GetString("kafka.host")
	addr := strings.Split(hosts, ",")

	kafka.P = kafka.NewProducer(addr)
}
