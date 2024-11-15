package challenge

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type Challange struct {
	Que      string
	Hash     string
	Original string
}

func randstring(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read")
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func hash(a string) string {
	h := sha256.Sum256([]byte(a))
	return base64.URLEncoding.EncodeToString(h[:])
}

func mask(s string, n int) string {
	q := []rune(s)
	for i := 0; i < n; i++ {
		q[i] = ' '
	}
	return string(q)
}

func newChallenge(n int) (Challange, error) {
	o, err := randstring(n)
	if err != nil {
		return Challange{}, fmt.Errorf("randstring: %w", err)
	}
	h := hash(o)
	q := mask(o, n)
	return Challange{Que: q, Hash: h, Original: o}, nil
}

func NewBatch(n, m int) ([]Challange, error) {
	cs := []Challange{}
	for range m {
		c, err := newChallenge(n)
		if err != nil {
			return nil, fmt.Errorf("newchallenge: %w", err)
		}
		cs = append(cs, c)
	}
	return cs, nil
}
