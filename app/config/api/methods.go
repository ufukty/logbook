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

type Child interface {
	GetParent() any
}

// endpoint/../.. = service
// service/../.. = gateway
func twoup(a Addressable) Addressable {
	c1, ok := a.(Child)
	if !ok {
		return nil
	}
	p1 := c1.GetParent()
	c2, ok := p1.(Child)
	if !ok {
		return nil
	}
	p2 := c2.GetParent()
	u, ok := p2.(Addressable)
	if !ok {
		return nil
	}
	return u
}

func PathFromInternet(a Addressable) string {
	addressables := []Addressable{a}
	for i := 0; i < 2; i++ { // endpoint > service (../) > gateway (../../)
		up := twoup(addressables[len(addressables)-1])
		if up == nil {
			break
		}
		addressables = append(addressables, up)
	}
	slices.Reverse(addressables)
	return Join(addressables...)
}
