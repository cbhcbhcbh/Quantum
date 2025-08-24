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
	Get(ctx context.Context, key string, dst any) (bool, error)
	Set(ctx context.Context, key string, val any) error
	Delete(ctx context.Context, key string) error
	HSet(ctx context.Context, key string, values ...any) error
	HGetIfKeyExists(ctx context.Context, key, field string, dst any) (bool, bool, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key, field string) error
	ExecPipeLine(ctx context.Context, cmds *[]RedisCmd) error
}

type RedisCacheImpl struct {
	client redis.UniversalClient
}

type RedisOpType int

const (
	DELETE RedisOpType = iota
	HSETONE
	RPUSH
)

type RedisPayload interface {
	Payload()
}

type RedisDeletePayload struct {
	RedisPayload
	Key string
}

type RedisHsetOnePayload struct {
	RedisPayload
	Key   string
	Field string
	Val   any
}

type RedisRpushPayload struct {
	RedisPayload
	Key string
	Val any
}

func (RedisDeletePayload) Payload()  {}
func (RedisHsetOnePayload) Payload() {}
func (RedisRpushPayload) Payload()   {}

type RedisCmd struct {
	OpType  RedisOpType
	Payload RedisPayload
}

type RedisPipelineCmd struct {
	OpType RedisOpType
	Cmd    any
}

func NewRedisClient(config *config.Config) (redis.UniversalClient, error) {
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

func (rc *RedisCacheImpl) Get(ctx context.Context, key string, dst any) (bool, error) {
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

func (rc *RedisCacheImpl) Set(ctx context.Context, key string, val any) error {
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

func (rc *RedisCacheImpl) HSet(ctx context.Context, key string, values ...any) error {
	return rc.client.HSet(ctx, key, values).Err()
}

var hgetIfKeyExists = redis.NewScript(`
local key = KEYS[1]
local field = ARGV[1]

if redis.call("EXISTS", key) == 0 then
  return ""
end

return redis.call("HGET", key, field)
`)

func (rc *RedisCacheImpl) HGetIfKeyExists(ctx context.Context, key, field string, dst any) (bool, bool, error) {
	val, err := hgetIfKeyExists.Run(ctx, rc.client, []string{key}, field).Text()
	if err == redis.Nil {
		return true, false, nil
	} else if err != nil {
		return false, false, err
	} else if val == "" {
		return false, false, nil
	} else {
		if err = json.Unmarshal([]byte(val), dst); err != nil {
			return false, false, err
		}
	}
	return true, true, nil
}

func (rc *RedisCacheImpl) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return rc.client.HGetAll(ctx, key).Result()
}

func (rc *RedisCacheImpl) HDel(ctx context.Context, key, field string) error {
	return rc.client.HDel(ctx, key, field).Err()
}

func (rc *RedisCacheImpl) ExecPipeLine(ctx context.Context, cmds *[]RedisCmd) error {
	pipe := rc.client.Pipeline()
	var pipelineCmds []RedisPipelineCmd

	for _, cmd := range *cmds {
		switch cmd.OpType {
		case DELETE:
			pipelineCmds = append(pipelineCmds, RedisPipelineCmd{
				OpType: DELETE,
				Cmd:    pipe.Del(ctx, cmd.Payload.(RedisDeletePayload).Key),
			})
		case HSETONE:
			payload := cmd.Payload.(RedisHsetOnePayload)
			pipelineCmds = append(pipelineCmds, RedisPipelineCmd{
				OpType: HSETONE,
				Cmd:    pipe.HSet(ctx, payload.Key, payload.Field, payload.Val),
			})
		case RPUSH:
			payload := cmd.Payload.(RedisRpushPayload)
			pipelineCmds = append(pipelineCmds, RedisPipelineCmd{
				OpType: RPUSH,
				Cmd:    pipe.RPush(ctx, payload.Key, payload.Val),
			})
		default:
			return ErrRedisPipelineCmdNotFound
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	for _, executedCmd := range pipelineCmds {
		switch executedCmd.OpType {
		case DELETE:
			if err := executedCmd.Cmd.(*redis.IntCmd).Err(); err != nil {
				return err
			}
		case HSETONE:
			if err := executedCmd.Cmd.(*redis.IntCmd).Err(); err != nil {
				return err
			}
		case RPUSH:
			if err := executedCmd.Cmd.(*redis.IntCmd).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}
