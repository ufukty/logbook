package endpoints

import (
	"errors"
)

var generic = "Unable to process your request at the moment"

var specific = map[error]string{}

func redact(err error) string {
	str := generic
	for ; err != nil; err = errors.Unwrap(err) {
		if s, ok := specific[err]; ok {
			str = s
			break
		}
	}
	return str
}
