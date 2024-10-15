package api

import (
	"path/filepath"
	"slices"
)

type Addressable interface {
	GetPath() string
}

func Join(addrs ...Addressable) string {
	j := ""
	for _, s := range addrs {
		j = filepath.Join(j, s.GetPath())
	}
	return j
}

type backtraceable interface {
	GetParent() any
}

// follows .Parent refs until finds an [addressable]
func up(a Addressable) Addressable {
	var cursor any = a
	for {
		btl, ok := cursor.(backtraceable)
		if !ok {
			return nil
		}
		cursor = btl.GetParent()
		if cursor == nil {
			return nil
		}
		a2, ok := cursor.(Addressable)
		if ok {
			return a2
		}
	}
}

func ancestry(a Addressable) []Addressable {
	ads := []Addressable{}
	for cursor := a; cursor != nil; cursor = up(cursor) {
		ads = append(ads, cursor)
	}
	slices.Reverse(ads)
	return ads
}

// returns service[,endpoint]
func ByService(a Addressable) string {
	return Join(ancestry(a)[1:]...)
}

// returns gateway[,service[,endpoint]]
func ByGateway(a Addressable) string {
	return Join(ancestry(a)...)
}

// returns gateway[,service[,endpoint]]
func PrefixedByGateway(a Addressable) string {
	return Join(ancestry(a)...) + "/"
}
