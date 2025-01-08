package cacher

type Cacher[T any] interface {
	Get(key string) (T, error)
	Set(key string, value any) error
	Del(key string) error
	Key() string
}
