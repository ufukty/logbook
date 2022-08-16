package crypto

import (
	"bytes"
	"testing"
)

func TestNewSalt(t *testing.T) {
	for i := 0; i < 1000; i++ {
		saltInBytes, err := NewSalt()
		if err != nil {
			t.Fail()
		}
		// log.Println("saltInBytes: ", saltInBytes)

		saltEncoded := Base64Encode(saltInBytes)
		// log.Println("saltEncoded: ", saltEncoded)

		saltDecoded, err := Base64Decode(saltEncoded)
		if err != nil {
			// log.Println(errors.Wrap(err, "Decoding encoded salt is failed"))
			t.Fail()
		}
		// log.Println("saltDecoded: ", saltDecoded)

		if !bytes.Equal(saltDecoded, saltInBytes) {
			// log.Println("Encoded decoded version is different than original salt.")
			t.Fail()
		}
	}
}
