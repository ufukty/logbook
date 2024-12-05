package challenge

import (
	"crypto/rand"
	"fmt"
)

func randbytes(l int) ([]byte, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("read")
	}
	return b, nil
}

var ErrEmptyAlphabet = fmt.Errorf("empty alphabet")

func randstring(alphabet string, l int) (string, error) {
	if len(alphabet) == 0 {
		return "", ErrEmptyAlphabet
	}
	bs, err := randbytes(l)
	if err != nil {
		return "", fmt.Errorf("randbytes: %w", err)
	}
	s := make([]byte, l)
	for i, b := range bs {
		s[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(s), nil
}
