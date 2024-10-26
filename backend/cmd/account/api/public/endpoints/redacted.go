package endpoints

import (
	"errors"
	"logbook/cmd/account/api/public/app"
	"logbook/cmd/account/permissions"
	"logbook/cmd/account/sessions"
)

var generic = "Unable to process your request at the moment"

var redacted = map[error]string{
	app.ErrEmailExists:           "Unable to process your request",
	app.ErrExpiredSession:        "Session has expired",
	app.ErrHashMismatch:          "Either the account doesn't exist or the credentials don't match",
	app.ErrSessionNotFound:       "Login is required",
	permissions.ErrUnauthorized:  "Higher authorization is required",
	sessions.ErrNoAuthentication: "Login is required",
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
