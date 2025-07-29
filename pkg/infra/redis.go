package infra

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/cbhcbhcbh/Quantum/pkg/common/util"
	"github.com/cbhcbhcbh/Quantum/pkg/config"
	redis "github.com/redis/go-redis/v9"
)

var (
	RedisClient                 redis.UniversalClient
	ErrRedisUnlockFail          = errors.New("redis unlock fail")
	ErrRedisPipelineCmdNotFound = errors.New("redis pipeline command not found; supports only SET and DELETE")

	expiration time.Duration
)

type RedisCache interface {
	Get(ctx context.Context, key string, dst interface{}) (bool, error)
	Set(ctx context.Context, key string, val interface{}) error
	Delete(ctx context.Context, key string) error
}

type RedisCacheImpl struct {
	client redis.UniversalClient
}

func NewredisClient(config *config.Config) (redis.UniversalClient, error) {
	expiration = time.Duration(config.Redis.ExpirationHour) * time.Hour
	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          util.GetServerAddrs(config.Redis.Addrs),
		Password:       config.Redis.Password,
		ReadOnly:       true,
		RouteByLatency: true,
		MinIdleConns:   config.Redis.MinIdleConn,
		PoolSize:       config.Redis.PoolSize,
		ReadTimeout:    time.Duration(config.Redis.ReadTimeoutMilliSecond) * time.Millisecond,
		WriteTimeout:   time.Duration(config.Redis.WriteTimeoutMilliSecond) * time.Millisecond,
		PoolTimeout:    5 * time.Second,
	})
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err == redis.Nil || err != nil {
		return nil, err
	}
	return RedisClient, nil
}

func NewRedisCacheImpl(client redis.UniversalClient) *RedisCacheImpl {
	return &RedisCacheImpl{client: client}
}

func (rc *RedisCacheImpl) Get(ctx context.Context, key string, dst interface{}) (bool, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		if err = json.Unmarshal([]byte(val), dst); err != nil {
			return false, err
		}
	}
	return true, nil
}

func (rc *RedisCacheImpl) Set(ctx context.Context, key string, val interface{}) error {
	if err := rc.client.Set(ctx, key, val, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *RedisCacheImpl) Delete(ctx context.Context, key string) error {
	if err := rc.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
