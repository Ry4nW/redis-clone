package cache

import (
	"redis-clone/internal/resp"
	"sync"
)

// aliasing
type RespValue = resp.RespValue

type GenericMapData[K comparable, V any] struct {
	data map[K]V
}

type Store struct {
	// shared map
	mu   sync.RWMutex
	data map[string]RespValue
}

func newStore() *Store {
	return &Store{
		data: make(map[string]RespValue),
	}
}

func (store *Store) Get(key string) (RespValue, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	res, ok := store.data[key]
	return res, ok
}

func (store *Store) Set(key string, val RespValue) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	store.data[key] = val
}
