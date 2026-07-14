package command

// import (
// 	"redis-clone/internal/resp"
// )

// // aliasing
// type RespValue = resp.RespValue

// type GenericMapData[K comparable, V any] struct {
// 	data map[K]V
// }

// type Store struct {
// 	store map[string]RespValue
// }

// func newStore() *Store {
// 	return &Store{
// 		store: make(map[string]RespValue),
// 	}
// }

// func (store *Store) Get(key string) (RespValue, bool) {
// 	res, ok := store.store[key]
// 	return res, ok
// }

// func (store *Store) Set(key string, val RespValue) {
// 	store.store[key] = val
// }
