package challenge

import (
	"fmt"
	"slices"
)

func start(n int) prefix {
	return prefix(slices.Repeat([]byte{alphabet[0]}, n))
}

type prefix string

func (p *prefix) iterate() bool {
	pb := []byte(*p)
	for i := len(pb) - 1; i >= 0; i-- {
		if i == 0 && pb[i] == alphabet[len(alphabet)-1] {
			return false
		}
		j := slices.Index([]byte(alphabet), pb[i])
		pb[i] = alphabet[(j+1)%len(alphabet)]
		*p = prefix(string(pb))
		if pb[i] != alphabet[0] {
			return true
		}
	}
	return true
}

func Solve(n int, que, hash_ string) (string, error) {
	if len(hash_) == 0 || n == 0 {
		return "", fmt.Errorf("invalid challenge: que or hash is empty")
	}
	for p := start(n); p.iterate(); {
		cand := string(p) + que
		if hash(cand) == hash_ {
			return cand, nil
		}
	}
	return "", fmt.Errorf("not found")
}
