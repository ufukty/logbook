package main

import (
	"fmt"
	"logbook/cmd/tasks/controllers/document"
	"logbook/cmd/tasks/controllers/task"
	"net/http"

	"github.com/gorilla/mux"
)

type Endpoint struct {
	Path    string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
}

func rootURIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world. Read our documentation to start.")
}

func registerEndpoints(r *mux.Router) {
	var (
		documentOverviewHierarchical  = document.CDocumentOverviewHierarchical{}
		documentOverviewChronological = document.CDocumentOverviewChronological{}
	)

	r.HandleFunc("/", rootURIHandler).Methods("GET")

	r.HandleFunc("/document", document.Create).Methods("POST")
	r.HandleFunc("/document/overview/{document_id}", documentOverviewHierarchical.Responder).Methods("GET")
	r.HandleFunc("/document/overview/chronological/{document_id}", documentOverviewChronological.Responder).Methods("GET").Queries("limit", "{limit}").Queries("offset", "{offset}")

	r.HandleFunc("/task", task.Create).Methods("POST")

	// {"/group/{document_id}", "GET", group.Get},

	// {"/task/{document_id}/{task_id}", "GET", task.Read},
	// {"/task/{document_id}/{task_id}", "PATCH", task.Update},
	// {"/task/{document_id}/{task_id}", "DELETE", task.Delete},
}
