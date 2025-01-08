package set

import (
	"context"
	"github.com/sosumecho/modules/drivers/redis"

	"time"
)

type RedisCache struct {
	ctx       context.Context
	ExpiredAt time.Duration
}

func (r RedisCache) Get(key string) (string, error) {
	return redis.Redis.SRandMember(r.ctx, key).Result()
}

func (r RedisCache) Set(key string, value string) error {
	return redis.Redis.SAdd(r.ctx, key, value).Err()
}

func (r RedisCache) Del(key string) error {
	if redis.Redis.Exists(r.ctx, key).Val() > 0 {
		return redis.Redis.Del(r.ctx, key).Err()
	}
	return nil
}

func (r RedisCache) DelItem(key string, item string) error {
	if redis.Redis.SIsMember(r.ctx, key, item).Val() {
		return redis.Redis.SRem(r.ctx, key, item).Err()
	}
	return nil
}

func (r RedisCache) Key() string {
	return ""
}

func NewRedisCache(expiredAt time.Duration) *RedisCache {
	return &RedisCache{
		ctx:       context.Background(),
		ExpiredAt: expiredAt,
	}
}
