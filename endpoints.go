package main

import (
	"fmt"
	"logbook/main/controller/document"
	"logbook/main/controller/task"

	"net/http"

	"github.com/gorilla/mux"
)

type Endpoint struct {
	Path    string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
}

func rootURIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logbook API is for 1st party use only.")
}

func registerEndpoints(r *mux.Router) {
	var (
		documentOverviewHierarchical  = document.CDocumentOverviewHierarchical{}
		documentOverviewChronological = document.CDocumentOverviewChronological{}
	)

	r.HandleFunc("/", rootURIHandler).Methods("GET")

	// r.HandleFunc("/auth/user", authUserCreate).
	// 	Methods("POST").
	// 	Queries("email_address", "{email_address}")

	// r.HandleFunc("/auth/session", authSessionCreate).
	// 	Methods("POST").
	// 	Queries("email_address", "{email_address}").
	// 	Queries("password_kdf", "{password_kdf}")

	r.HandleFunc("/document", document.Create).Methods("POST")
	r.HandleFunc("/document/{id}/placement/hierarchical", documentOverviewHierarchical.Responder).Methods("GET")
	r.HandleFunc("/document/{id}/placement/chronological", documentOverviewChronological.Responder).Methods("GET").Queries("limit", "{limit}", "offset", "{offset}")

	r.HandleFunc("/task", task.Create).Methods("POST")

}
