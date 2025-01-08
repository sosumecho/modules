package pool

import (
	"github.com/panjf2000/ants/v2"
	"sync"
)

var (
	pool *ants.Pool
	once sync.Once
)

func Pool() *ants.Pool {
	if pool == nil {
		once.Do(func() {
			var err error
			pool, err = ants.NewPool(100)
			if err != nil {
				panic(err)
			}
		})
	}
	return pool
}

func NewPool(size int) *ants.Pool {
	p, err := ants.NewPool(size, ants.WithNonblocking(false))
	if err != nil {
		return nil
	}
	return p
}
