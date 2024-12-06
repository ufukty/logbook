package challenge

import (
	"fmt"
	"slices"
)

func start(difficulty int) prefix {
	return prefix{
		value:      slices.Repeat([]byte{alphabet[0]}, ML),
		difficulty: difficulty,
	}
}

type prefix struct {
	value      []byte
	difficulty int
}

func (p *prefix) iterate() bool {
	pb := slices.Clone(p.value)
	for i := len(pb) - 1; i >= 0; i-- {
		if i == 0 && pb[i] == alphabet[p.difficulty-1] {
			return false
		}
		j := slices.Index([]byte(alphabet), pb[i])
		pb[i] = alphabet[(j+1)%p.difficulty]
		p.value = pb
		if pb[i] != alphabet[0] {
			return true
		}
	}
	return true
}

var ErrNotFound = fmt.Errorf("not found")

func SolveChallenge(difficulty int, masked, hashed string) (string, error) {
	if difficulty <= 2 {
		return "", ErrMinDifficulty
	}
	if len(alphabet) <= difficulty {
		return "", ErrMaxDifficulty
	}
	if len(hashed) == 0 {
		return "", fmt.Errorf("hashed is empty")
	}
	if len(masked) == 0 {
		return "", fmt.Errorf("masked is empty")
	}
	p := start(difficulty)
	for {
		cand := string(p.value) + masked
		if hash(cand) == hashed {
			return string(p.value), nil
		}
		if !p.iterate() {
			return "", ErrNotFound
		}
	}
}
