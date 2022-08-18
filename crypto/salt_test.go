package crypto

import (
	"bytes"
	"sort"
	"testing"
)

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func IsEqual(s1, s2 []byte) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func TestNewSaltEncodingDecoding(t *testing.T) {
	for i := 0; i < 100000; i++ {
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

func TestSaltCollusions(t *testing.T) {

	producedSalts := [][]byte{}

	for i := 0; i < 10000000; i++ {

		saltInBytes, err := NewSalt()
		if err != nil {
			t.Fail()
		}
		producedSalts = append(producedSalts, saltInBytes)
	}

	sort.Slice(producedSalts, func(i, j int) bool {
		// returns true when i < j
		v1 := producedSalts[i]
		v2 := producedSalts[i]
		for p := 0; p < Min(len(v1), len(v2)); p++ {
			if v1[p] < v2[p] {
				return true
			} else if v1[p] > v2[p] {
				return false
			}
			// if they are same iterate to next p
		}
		return len(v1) < len(v2)
	})

	collusions := 0
	for i := 0; i < len(producedSalts)-1; i++ {
		if IsEqual(producedSalts[i], producedSalts[i+1]) {
			collusions += 1
		}
	}

	if collusions > 0 {
		t.Errorf("More than one collusions (%d) on produced salts.", collusions)
	}

}
