package challenge

import (
	"crypto/sha256"
)

func hash(a string) string {
	h := sha256.Sum256([]byte(a))
	return encode(h[:])
}
