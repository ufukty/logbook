package app

import "sync"

type Set[T comparable] struct {
	mu    sync.RWMutex
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

func (s *Set[T]) Add(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item] = struct{}{}
}

func (s *Set[T]) Delete(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, item)
}

func (s *Set[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.items[item]
	return exists
}

func (s *Set[T]) Items() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]T, 0, len(s.items))
	for item := range s.items {
		items = append(items, item)
	}
	return items
}
