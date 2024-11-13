package endpoints

import (
	"errors"
	"logbook/cmd/sessions/app"
)

var generic = "Unable to process your request at the moment"

var redacted = map[error]string{
	app.ErrExpiredSession:  "Session has expired",
	app.ErrHashMismatch:    "Either the account doesn't exist or the credentials don't match",
	app.ErrSessionNotFound: "Login is required",
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
