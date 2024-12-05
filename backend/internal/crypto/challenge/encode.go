package challenge

import "encoding/base64"

var encoder = base64.StdEncoding.WithPadding(base64.NoPadding)

func encode(s []byte) string {
	return encoder.EncodeToString(s)
}
