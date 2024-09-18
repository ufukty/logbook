package stores

type Stores[K comparable, V any] interface {
	Len(k K, v V) int
	Has(k K) bool
	Set(k K, v V)
	Get(k K) (V, bool)
	Delete(k K)
	Keys() []K
}
