package transports

import (
	"crypto/rand"
	"encoding/base64"
)

func NewAntiCsrfToken() (AntiCsrfToken, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "reading random bytes", err
	}
	return AntiCsrfToken(base64.URLEncoding.EncodeToString(b)), nil
}
