package gsync

import "sync"

// Wrapper of [sync.Map] supports generics. See [sync.Map] for docs
type Map[K comparable, V any] struct {
	m sync.Map
}

func (m *Map[K, V]) Clear() {
	m.m.Clear()
}

func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.m.Load(key)
	return v.(V), ok
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, l := m.m.LoadAndDelete(key)
	return v.(V), l
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, l := m.m.LoadOrStore(key, value)
	return v.(V), l
}

func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	p, l := m.m.Swap(key, value)
	return p.(V), l
}
