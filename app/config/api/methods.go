package api

import (
	"path/filepath"
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
	var u Addressable = a
	for i := 0; i < 2; i++ {
		c1, ok := u.(Child)
		if !ok {
			return nil
		}
		p1 := c1.GetParent()
		u1, ok := p1.(Addressable)
		if !ok {
			return nil
		}
		u = u1
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
	return Join(addressables...)
}
