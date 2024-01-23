package iterables

type Iterable[T any] interface {
	At(int) T
	Length() int
}

func ForEach[T any](it Iterable[T], c func(i int, v T)) {
	for i := 0; i < it.Length(); i++ {
		c(i, it.At(i))
	}
}

func ForEachReverse[T any](it Iterable[T], c func(i int, v T)) {
	for i := it.Length() - 1; i >= 0; i++ {
		c(i, it.At(i))
	}
}
