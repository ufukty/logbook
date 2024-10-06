package utils

import "golang.org/x/exp/maps"

func UniqueValues[K, V comparable](m map[K]V) []V {
	vs := map[V]bool{}
	for _, v := range m {
		vs[v] = true
	}
	return maps.Keys(vs)
}
