package crypto

import (
	"log"
	"testing"

	"github.com/pkg/errors"
)

func TestBase64(t *testing.T) {

	testCases := []string{
		"Hello world",    // SGVsbG8gd29ybGQ=
		"Hello world.",   // SGVsbG8gd29ybGQu
		"Hello worlds.",  // SGVsbG8gd29ybGRzLg==
		"Hello world...", // SGVsbG8gd29ybGQuLi4=
	}

	for _, text := range testCases {
		log.Println("text: ", text)

		encoded := Base64Encode([]byte(text))
		log.Println("encoded: ", encoded, " string-representation: ", string(encoded))

		encodedStripped := StripEndPadding(encoded)
		log.Println("encodedStripped: ", encodedStripped, " string-representation: ", string(encodedStripped))

		encodedPadded := AddPadding(encodedStripped)
		log.Println("encodedPadded: ", encodedPadded, " string-representation: ", string(encodedPadded))

		decoded, err := Base64Decode(encodedPadded)
		if err != nil {
			log.Println(errors.Wrap(err, "TestBase64"))
			t.Fail()
		}
		log.Println("decoded: ", decoded)

		decodedString := string(decoded)
		if decodedString != text {
			log.Println("Encoded/Stripped/Padded/Decoded version of text is not same with original.")
			t.Fail()
		}
		log.Println("decodedString: ", decodedString)
	}

	faultyDecodingCases := []string{
		"SGVsbG8gd29ybGRzLg========",
		"SGVsbG8gd29ybGRzLg=======",
		"SGVsbG8gd29ybGRzLg=====",
		"SGVsbG8gd29ybGRzLg====",
		"SGVsbG8gd29ybGRzLg===",
	}

	for _, text := range faultyDecodingCases {
		_, err := Base64Decode([]byte(text))
		if err == nil {
			log.Println(errors.Wrap(err, "TestBase64 is accepted faulty input and didn't return error"))
			t.Fail()
		}
		log.Println(errors.Wrap(err, "Base64Decode returned error on faulty input"))
	}
}
