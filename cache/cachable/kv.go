package cachable

import (
	"github.com/sosumecho/modules/cache/cacher"
	"github.com/sosumecho/modules/cache/encoder"
	"github.com/sosumecho/modules/cache/kv"
	"time"
)

type KvCache[T any] struct {
	key string
	cacher.Cacher[T]
}

func (k KvCache[T]) Key() string {
	return k.key
}

func (k KvCache[T]) GetOrForget(key string) (T, error) {
	rs, err := k.Get(key)
	k.Del(key)
	return rs, err
}

func NewKV[T any](key string, expired time.Duration) *KvCache[T] {
	jsonEncoder := encoder.New[T]("json")
	c := kv.NewRedisCache(jsonEncoder, expired)
	return &KvCache[T]{
		key:    key,
		Cacher: c,
	}
}
