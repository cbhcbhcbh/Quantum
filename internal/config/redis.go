package config

import (
	"time"

	goRedis "github.com/cbhcbhcbh/Quantum/pkg/redis"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var RedisDB *redis.Client

func InitRedis() {
	redisOptions := &redis.Options{
		Addr:     viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),

		PoolSize:     viper.GetInt("redis.poll"),
		MinIdleConns: viper.GetInt("redis.conn"),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolTimeout:  5 * time.Second,
	}

	RedisDB = goRedis.NewRedisClient(redisOptions)
}
