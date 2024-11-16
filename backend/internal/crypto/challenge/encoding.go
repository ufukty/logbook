package challenge

import (
	"encoding/base32"
)

var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

var encoder = base32.StdEncoding.WithPadding(base32.NoPadding)

func encode(s []byte) string {
	return encoder.EncodeToString(s)
}
