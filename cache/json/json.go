package json

import (
	"github.com/sosumecho/modules/cache/cache"
	"github.com/sosumecho/modules/cache/cacher"
	"github.com/sosumecho/modules/cache/encoder"
	"github.com/sosumecho/modules/cache/kv"
	"time"
)

type Cache[T any] struct {
	cacher.Cacher[T]
	key string
}

func (c Cache[T]) Key() string {
	return c.key
}

func NewCache[T any](key string, expire time.Duration) *Cache[T] {
	jsonEncoder := encoder.New[T]("json")
	c := kv.NewRedisCache[T](jsonEncoder, expire)
	return &Cache[T]{
		Cacher: c,
		key:    key,
	}
}

func Json[T any](key string, expire time.Duration, cacheAble cache.Callback[T]) (T, error) {
	return cache.Cache[T](NewCache[T](key, expire), cacheAble)
}
