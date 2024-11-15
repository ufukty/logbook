package challenge

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func Solve(que, hash string) (string, error) {
	if len(hash) == 0 || len(que) == 0 {
		return "", fmt.Errorf("invalid challenge: que or hash is empty")
	}

	for attempt := 0; true; attempt++ {
		fmt.Println("attempt:", attempt)
		candidate, err := generateCandidate(len(que), attempt)
		if err != nil {
			return "", fmt.Errorf("failed to generate candidate: %w", err)
		}

		h := sha256.Sum256([]byte(candidate))
		cand := base64.URLEncoding.EncodeToString(h[:])
		if cand == hash {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("not found")
}

func generateCandidate(length, attempt int) (string, error) {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte((attempt + i) % 256)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
