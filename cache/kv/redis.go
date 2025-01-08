package kv

import (
	"context"
	"github.com/sosumecho/modules/cache/cacher"
	"github.com/sosumecho/modules/cache/encoder"
	"github.com/sosumecho/modules/drivers/redis"

	"time"
)

type RedisCache[T any] struct {
	ctx       context.Context
	encoder   encoder.Encoder[T]
	expiredAt time.Duration
}

func (r RedisCache[T]) Get(key string) (T, error) {
	var (
		rs     T
		result string
		err    error
	)

	result, err = redis.Redis.Get(r.ctx, key).Result()
	if err != nil {
		return rs, err
	}
	rs, err = r.encoder.Decode([]byte(result))
	if err != nil {
		return rs, err
	}
	return rs, nil

}

func (r RedisCache[T]) Set(key string, value any) error {
	rs, err := r.encoder.Encode(value)
	if err != nil {
		return err
	}
	return redis.Redis.Set(r.ctx, key, string(rs), r.expiredAt).Err()
}

func (r RedisCache[T]) Del(key string) error {
	if redis.Redis.Exists(r.ctx, key).Val() > 0 {
		return redis.Redis.Del(r.ctx, key).Err()
	}
	return nil
}

func (r RedisCache[T]) Key() string {
	return ""
}

func NewRedisCache[T any](encoder encoder.Encoder[T], expiredAt time.Duration) cacher.Cacher[T] {
	return RedisCache[T]{
		ctx:       context.Background(),
		expiredAt: expiredAt,
		encoder:   encoder,
	}
}
