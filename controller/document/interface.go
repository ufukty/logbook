package document

import (
	e "logbook/main/controller/utilities/errors"
	db "logbook/main/database"
	"net/http"
)

type IController interface {
	sanitizer(r *http.Request) (string, []error)
	executor(r *http.Request) ([]db.Task, *e.Error)
	Responder(w http.ResponseWriter, r *http.Request)
}