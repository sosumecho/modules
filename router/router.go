package router

import (
	"fmt"
	"strings"
	"sync"
)

var (
	manager *RouteManager
	once    sync.Once
)

type RouteManager struct {
	lock  sync.RWMutex
	Alias map[string]string
}

func NewRouteManager() *RouteManager {
	if manager == nil {
		once.Do(func() {
			manager = &RouteManager{
				Alias: make(map[string]string),
			}
		})
	}
	return manager
}

func (r *RouteManager) Key(path, method string) string {
	path = strings.TrimPrefix(path, "/admin")
	path = strings.TrimPrefix(path, "/api")
	return fmt.Sprintf("%s:%s", method, path)
}

func (r *RouteManager) Add(path, method, alias string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.Alias[r.Key(path, method)] = alias
}

func (r *RouteManager) Get(path, method string) (string, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	val, ok := r.Alias[r.Key(path, method)]
	return val, ok
}
