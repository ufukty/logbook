package document

import (
	"logbook/cmd/tasks/database"
	e "logbook/controller/utilities/errors"
	"net/http"
)

type IController interface {
	sanitizer(r *http.Request) (string, []error)
	executor(r *http.Request) ([]database.Task, *e.Error)
	Responder(w http.ResponseWriter, r *http.Request)
}
