package api

import (
	"path/filepath"
	"slices"
)

type addressable interface {
	GetPath() string
}

func Join(addrs ...addressable) string {
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
func up(a addressable) addressable {
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
		a2, ok := cursor.(addressable)
		if ok {
			return a2
		}
	}
}

func ancestry(a addressable) []addressable {
	ads := []addressable{}
	for cursor := a; cursor != nil; cursor = up(cursor) {
		ads = append(ads, cursor)
	}
	slices.Reverse(ads)
	return ads
}

// returns service[,endpoint]
func ByService(a addressable) string {
	return Join(ancestry(a)[1:]...)
}

// returns gateway[,service[,endpoint]]
func ByGateway(a addressable) string {
	return Join(ancestry(a)...)
}

// returns gateway[,service[,endpoint]]
func PrefixedByGateway(a addressable) string {
	return Join(ancestry(a)...) + "/"
}
