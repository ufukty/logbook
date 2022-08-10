package controllers

import (
	"net/http"

	parameters "logbook/main/parameters"
	responder "logbook/main/responder"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	params := parameters.UserCreate{}
	if err := params.InputSanitizer(r); err != nil {
		responder.ErrorHandler(w, r, http.StatusBadRequest, "Check your parameters", err)
	}

	// check CSRF

	responder.SuccessHandler(w, r, params.Response)
}


func TaskCreate(w http.ResponseWriter, r *http.Request) {
	params := parameters.TaskCreate{}
	if err := params.InputSanitizer(r); err != nil {
		responder.ErrorHandler(w, r, http.StatusBadRequest, "Check your parameters", err)
	}

	// check auth

	// check existence of super task

	// check permissions between task and user

	// create NewOperation

	responder.SuccessHandler(w, r, params.Response)
}