package stores

import (
	"logbook/internal/logger"
	"sync"
)

type FixedSizeKV[K comparable, V any] struct {
	Size  int
	store map[K]V
	mu    sync.RWMutex
	log   logger.Logger
}

var _ Stores[any, any] = &FixedSizeKV[any, any]{}

func NewFixedSizeKV[K comparable, V any](logname string, size int) *FixedSizeKV[K, V] {
	return &FixedSizeKV[K, V]{
		Size:  size,
		store: map[K]V{},
		mu:    sync.RWMutex{},
		log:   *logger.New(logname),
	}
}

// TODO: implement LRU cache, if it'll perform better on storing integers (as the timestamps will consume memory too)
func (f *FixedSizeKV[K, V]) maintainsize() {
	if len(f.store) < f.Size {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.log.Println("reducing size")
	keys := make([]K, 0, len(f.store))
	for key := range f.store {
		keys = append(keys, key)
	}
	for _, key := range keys[:int(f.Size/2)] {
		delete(f.store, key)
	}
}

func (f *FixedSizeKV[K, V]) Len(k K, v V) int {
	return len(f.store)
}

func (f *FixedSizeKV[K, V]) Has(k K) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.store[k]
	return ok
}

func (f *FixedSizeKV[K, V]) Set(k K, v V) {
	f.maintainsize()
	f.mu.Lock()
	defer f.mu.Unlock()
	f.store[k] = v
}

func (f *FixedSizeKV[K, V]) Get(k K) (V, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	v, ok := f.store[k]
	return v, ok
}

func (f *FixedSizeKV[K, V]) Delete(k K) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.store, k)
}

func (f *FixedSizeKV[K, V]) Keys() []K {
	f.mu.RLock()
	defer f.mu.RUnlock()
	keys := make([]K, 0, len(f.store))
	for key := range f.store {
		keys = append(keys, key)
	}
	return keys
}
