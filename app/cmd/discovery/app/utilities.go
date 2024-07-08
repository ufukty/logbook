package app

type Set[C comparable] map[C]bool

func (s *Set[C]) Delete(c C) {
	delete(*s, c)
}

func (s *Set[C]) Add(c C) {
	(*s)[c] = true
}
