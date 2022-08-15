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

func VerifyHash(storedHash, clearText string) (bool, error) {
	result, err := argon2.VerifyEncoded(storedHash, []byte(clearText))
	if err != nil {
		return false, errors.Wrap(err, "VerifyHash()")
	}
	return result, nil
}

func Hash(clear, salt string) (string, error) {
	hashedString, err := argon2.HashEncoded(defaultArgon2Config, []byte(clear), []byte(salt))
	if err != nil {
		return "", errors.Wrap(err, "Hash()")
	}
	return hashedString, nil
}
