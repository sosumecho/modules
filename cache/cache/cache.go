package cache

import (
	"github.com/sosumecho/modules/cache/cacher"
)

type Callback[T any] func() (T, error)

func Cache[T any](cacher cacher.Cacher[T], cacheAble Callback[T]) (T, error) {
	rs, err := cacher.Get(cacher.Key())
	if err != nil {
		rs, err = cacheAble()
		if err = cacher.Set(cacher.Key(), rs); err != nil {
			return rs, err
		}
	}
	return rs, err
}
