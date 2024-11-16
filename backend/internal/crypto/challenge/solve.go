package challenge

import (
	"fmt"
)

func Solve(n int, que, hash_ string) (string, error) {
	if len(hash_) == 0 || n == 0 {
		return "", fmt.Errorf("invalid challenge: que or hash is empty")
	}
	l := len(alphabet) ^ n
	for i := 0; i < l; i++ {
		comb := combinate(n, i)
		cand := comb + que
		if hash(cand) == hash_ {
			return cand, nil
		}
	}
	return "", fmt.Errorf("not found")
}

// generates the i'th value in the range from "AAA...A" to "7777...7"
func combinate(n int, i int) string {
	if n <= 0 || i < 0 {
		return ""
	}
	base := len(alphabet)
	result := make([]byte, n)
	for j := n - 1; j >= 0; j-- {
		result[j] = alphabet[i%base]
		i /= base
	}
	return string(result)
}
