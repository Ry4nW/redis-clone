package cache

import (
	"redis-clone/internal/resp"
)


// aliasing 
type RespValue = resp.RespValue


type GenericMapData[K comparable, V any] struct {
    data map[K]V
}


type Store struct {
	data map[string]RespValue
}

func newStore() * Store {
	return &Store{
		data: make(map[string]RespValue),
	}
}

func (store *Store) Get(key string) (RespValue,bool) {
	res, ok := store.data[key]
	return res, ok
}

func (store *Store) Set(key string, val RespValue)  {
	store.data[key] = val
}

