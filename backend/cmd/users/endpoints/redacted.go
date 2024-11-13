package endpoints

import (
	"errors"
)

var generic = "Unable to process your request at the moment"

var redacted = map[error]string{}

func redact(err error) string {
	s := generic
	for ; err != nil; err = errors.Unwrap(err) {
		if e, ok := redacted[err]; ok {
			s = e
			break
		}
	}
	return s
}
