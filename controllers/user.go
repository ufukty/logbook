package controllers

import (
	"logbook/main/parameters"
	"logbook/main/responder"
	"net/http"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	params := parameters.UserCreate{}
	if err := params.InputSanitizer(r); err != nil {
		responder.ErrorHandler(w, r, http.StatusBadRequest, "Check your parameters", err)
	}

	// check CSRF

	// create first task ➝ returning task_id

	// create first bookmark

	responder.SuccessHandler(w, r, params.Response)
}
