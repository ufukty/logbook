package crypto

import (
	"crypto/rand"
	"io"

	"github.com/pkg/errors"
)

const MAX_SALT_SIZE_IN_BYTES = 16

func NewSalt() ([]byte, error) {
	randomBytes := make([]byte, MAX_SALT_SIZE_IN_BYTES)
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return nil, errors.Wrap(err, "NewSalt()")
	}
	return randomBytes, nil
}
