package main

import (
	"logbook/main/controller/document"
	"logbook/main/controller/task"
	"net/http"
)

type Endpoint struct {
	Path    string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
}

func endpoints() []Endpoint {
	return []Endpoint{

		{"/document", "POST", document.Create},
		// {"/document/list", "GET", document.List},
		// {"/document/{document_id}", "GET", document.Details},
		// {"/document/{document_id}", "PATCH", document.Update},
		// {"/document/{document_id}", "DELETE", document.Delete},

		// {"/group/{document_id}", "GET", group.Get},

		// {"/task/{document_id}/{task_id}", "GET", task.Read},
		{"/task", "POST", task.Create},
		// {"/task/{document_id}/{task_id}", "PATCH", task.Update},
		// {"/task/{document_id}/{task_id}", "DELETE", task.Delete},
	}
}
