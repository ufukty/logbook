package crypto

import (
	"log"
	"testing"

	"github.com/pkg/errors"
	"github.com/tvdburgt/go-argon2"
)

func TestArgon2Hash(t *testing.T) {
	clearText := []byte("Hello world")
	salt := []byte("Lorem ipsum dolor sit amet")

	hash, err := Argon2Hash(clearText, salt)
	if err != nil {
		log.Println(errors.Wrap(err, "Hashing"))
		t.Fail()
	}
	log.Println("hash: ", hash)

	status, err := Argon2Verify(hash, clearText)
	if err != nil {
		log.Println(errors.Wrap(err, "Verification"))
		t.Fail()
	}
	if status != true {
		log.Println("Verification has failed for correct clearText")
		t.Fail()
	}

	status, err = Argon2Verify(hash, []byte("hello darkness"))
	if err != nil {
		log.Println(errors.Wrap(err, "Verification 2"))
		t.Fail()
	}
	if status != false {
		log.Println("Verification has succeeded for wrong clearText")
		t.Fail()
	}

	bytes, err := argon2.Hash(&argon2.Context{
		Iterations:  3,
		Memory:      1 << 12, // 4 MiB
		Parallelism: 1,
		HashLen:     32,
		Mode:        argon2.ModeArgon2id,
		Version:     argon2.Version13,
	}, clearText, salt)
	if err != nil {
		log.Println(errors.Wrap(err, "byyrd"))
	}
	log.Println(string(StripEndPadding(Base64Encode(salt))))
	log.Println(string(StripEndPadding(Base64Encode(bytes))))
}
