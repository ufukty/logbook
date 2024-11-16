package challenge

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

type Challange struct {
	Que      string
	Hash     string
	Original string
	N        int
}

func (c Challange) String() string {
	return fmt.Sprintf("(Que: %s) (Hash: %s) (Original: %s) (N: %d)\n", c.Que, c.Hash, c.Original, c.N)
}

func randstring(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read")
	}
	return encode(b), nil
}

func hash(a string) string {
	h := sha256.Sum256([]byte(a))
	return encode(h[:])
}

func mask(s string, n int) string {
	return string(s[n:])
}

func NewChallenge(l, n int) (Challange, error) {
	if l <= n {
		return Challange{}, fmt.Errorf("ques should be longer than masks")
	}
	o, err := randstring(l)
	if err != nil {
		return Challange{}, fmt.Errorf("randstring: %w", err)
	}
	h := hash(o)
	q := mask(o, n)
	c := Challange{
		N:        n,
		Que:      q,
		Hash:     h,
		Original: o,
	}
	return c, nil
}

func NewBatch(l, n, m int) ([]Challange, error) {
	cs := []Challange{}
	for range m {
		c, err := NewChallenge(l, n)
		if err != nil {
			return nil, fmt.Errorf("newchallenge: %w", err)
		}
		cs = append(cs, c)
	}
	return cs, nil
}
