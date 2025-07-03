package redis

import (
	"time"

	"github.com/cbhcbhcbh/Quantum/internal/pkg/known"
	"github.com/go-redis/redis"

	"sync"
)

var (
	once sync.Once
	R    *REDIS
)

type REDIS struct {
	redisClient *redis.Client
}

func NewredisClient(opts *redis.Options) *REDIS {

	once.Do(func() {
		R = &REDIS{
			redisClient: redis.NewClient(opts),
		}
	})

	return R
}

func (R *REDIS) SetKey(key string, value string, expiration time.Duration) {
	_ = R.redisClient.Set(key, value, expiration).Err()
}

func (R *REDIS) GetKey(key string) string {
	val, _ := R.redisClient.Get(key).Result()
	return val
}

func (R *REDIS) DelKey(key string) {
	_ = R.redisClient.Del(key).Err()
}

func (R *REDIS) SetBitmaps(id int64, status int) {
	_ = R.redisClient.SetBit(known.RedisBitmapUserLoggedKey, id, status)
}

func (R *REDIS) GetBitmaps(id int64) bool {
	status, _ := R.redisClient.GetBit(known.RedisBitmapUserLoggedKey, id).Result()
	return status == 1
}

func (R *REDIS) CountBitmaps() int64 {
	count, _ := R.redisClient.BitCount(known.RedisBitmapUserLoggedKey, nil).Result()
	return count
}
