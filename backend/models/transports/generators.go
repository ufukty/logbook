package transports

import (
	"crypto/rand"
	"encoding/base64"
)

func NewAntiCsrfToken() (AntiCsrfToken, error) {
	b := make([]byte, length_anti_csrf_token)
	if _, err := rand.Read(b); err != nil {
		return "reading random bytes", err
	}
	return AntiCsrfToken(base64.URLEncoding.EncodeToString(b)), nil
}
