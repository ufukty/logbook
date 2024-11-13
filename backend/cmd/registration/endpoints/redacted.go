package endpoints

import (
	"errors"
	"logbook/cmd/registration/app"
)

var generic = "Unable to process your request at the moment"

var redacted = map[error]string{
	app.ErrEmailExists: "Unable to process your request",
}

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
