package gsync

import (
	"sync"
)

// Wrapper of [sync.Pool] supports generics
type Pool[T any] struct {
	p sync.Pool
}

func (p *Pool[T]) Get() T {
	return p.p.Get().(T)
}
func (p *Pool[T]) Put(x T) {
	p.p.Put(x)
}
