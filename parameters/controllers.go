package parameters

import (
	"net/http"

	responder "logbook/main/responder"
)

func ControllerTaskCreate(w http.ResponseWriter, r *http.Request) {
	parameters := TaskCreate{}
	if err := parameters.InputSanitizer(r); err != nil {
		responder.ErrorHandler(w, r, http.StatusBadRequest, "Check your parameters", err)
	}

	// check auth

	// check existence of super task

	// check permissions between task and user

	responder.SuccessHandler(w, r, parameters.Response)
}

func ControllerUserCreate(w http.ResponseWriter, r *http.Request) {
	parameters := UserCreate{}
	if err := parameters.InputSanitizer(r); err != nil {
		responder.ErrorHandler(w, r, http.StatusBadRequest, "Check your parameters", err)
	}

	// check CSRF

	responder.SuccessHandler(w, r, parameters.Response)
}
