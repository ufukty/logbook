package crypto

import (
	"encoding/base64"
)

const base64Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var encoder = base64.NewEncoding(base64Alphabet)

func Base64Encode(input []byte) []byte {
	length := encoder.EncodedLen(len(input))
	output := make([]byte, length)
	encoder.Encode(output, input)
	return output
}

func StripEndPadding(input []byte) []byte {
	length := len(input)
	for input[length-1] == '=' {
		length--
	}
	return input[:length]
}

func AddPadding(input []byte) []byte {
	numOfEquals := (4 - len(input)%4) % 4
	for i := 0; i < numOfEquals; i++ {
		input = append(input, '=')
	}
	return input
}

func Base64Decode(input []byte) ([]byte, error) {
	output := make([]byte, encoder.DecodedLen(len(input)))
	n, err := encoder.Decode(output, input)
	if err != nil {
		return nil, err
	}
	return output[:n], nil
}
