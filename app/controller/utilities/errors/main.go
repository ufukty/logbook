package errors

import "net/http"

type Error struct {
	HttpResponseCode int
	HttpResponseHint string
	ErrorTrace       []error
}

func New(params ...interface{}) *Error {
	instance := Error{}
	// We need those just because GO doesn't support
	// optional parameters and default values
	httpResponseHintSet := false
	httpResponseCode := false
	for _, param := range params {
		switch v := param.(type) {
		case string:
			httpResponseHintSet = true
			instance.HttpResponseHint = v
		case int:
			httpResponseCode = true
			instance.HttpResponseCode = v
		case []error:
			instance.ErrorTrace = v
		}
	}
	if !httpResponseHintSet {
		instance.HttpResponseHint = "Check back later."
	}
	if !httpResponseCode {
		instance.HttpResponseCode = http.StatusInternalServerError
	}
	return &instance
}
