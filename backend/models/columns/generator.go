package columns

import (
	"fmt"

	"github.com/google/uuid"
)

func NewUuidV4[T ~string]() (T, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("uuid.NewRandom: %w", err)
	}
	return T(uuid.String()), nil
}

func NewUuidV4Unsafe[T ~string]() T {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("uuid.NewRandom: %w", err))
	}
	return T(uuid.String())
}
