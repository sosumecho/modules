package cacher

type SetCacher[T any] interface {
	Cacher[T]
	DelItem(key string, item string) error
}
