package stores

import (
	"sync"
)

type KV[K comparable, V any] struct {
	store map[K]V
	mu    sync.RWMutex
}

var _ Stores[any, any] = &KV[any, any]{}

func NewKV[K comparable, V any]() *KV[K, V] {
	return &KV[K, V]{
		store: map[K]V{},
		mu:    sync.RWMutex{},
	}
}

func (f *KV[K, V]) Len(k K, v V) int {
	return len(f.store)
}

func (f *KV[K, V]) Has(k K) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.store[k]
	return ok
}

func (f *KV[K, V]) Set(k K, v V) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.store[k] = v
}

func (f *KV[K, V]) Get(k K) (V, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	v, ok := f.store[k]
	return v, ok
}

func (f *KV[K, V]) Delete(k K) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.store, k)
}

func (f *KV[K, V]) Keys() []K {
	f.mu.RLock()
	defer f.mu.RUnlock()
	keys := make([]K, 0, len(f.store))
	for key := range f.store {
		keys = append(keys, key)
	}
	return keys
}
