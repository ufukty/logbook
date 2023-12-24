package crypto

import (
	"github.com/pkg/errors"
	"github.com/tvdburgt/go-argon2"
)

var defaultArgon2Config = &argon2.Context{
	Iterations:  3,
	Memory:      1 << 12, // 4 MiB
	Parallelism: 1,
	HashLen:     32,
	Mode:        argon2.ModeArgon2id,
	Version:     argon2.Version13,
}

func Argon2Hash(clearText []byte, salt []byte) (string, error) {
	hashedString, err := argon2.HashEncoded(defaultArgon2Config, clearText, salt)
	if err != nil {
		return "", errors.Wrap(err, "Argon2Hash()")
	}
	return hashedString, nil
}

func Argon2Verify(storedHash string, clearText []byte) (bool, error) {
	result, err := argon2.VerifyEncoded(storedHash, clearText)
	if err != nil {
		return false, errors.Wrap(err, "Argon2Verify()")
	}
	return result, nil
}
