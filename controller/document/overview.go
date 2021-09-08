package document

import (
	c "logbook/main/controller"
	e "logbook/main/controller/utilities/errors"
	v "logbook/main/controller/utilities/validate"
	db "logbook/main/database"
	"net/http"

	"github.com/gorilla/mux"
)

func sanitizeUserInput(r *http.Request) (string, []error) {
	vars := mux.Vars(r)
	documentId := vars["document_id"]
	if !v.IsValidUUID(documentId) {
		return "", []error{c.ErrDocumentIdInputIsNotValidUUID}
	}
	return documentId, nil
}

func overviewExecutor(r *http.Request) ([]db.Task, *e.Error) {
	var (
		tasks []db.Task
		errs  []error
	)

	documentId, errs := sanitizeUserInput(r)
	if errs != nil {
		return nil, e.New("Check your inputs.", http.StatusBadRequest, errs)
	}

	// create document table record
	tasks, errs = db.GetDocumentOverviewWithDocumentId(documentId)
	if errs != nil {
		return nil, e.New(errs)
	}

	return tasks, nil
}

func Overview(w http.ResponseWriter, r *http.Request) {
	// ipAddress := (*r).RemoteAddr
	// userAgent := (*r).Header.Get("User-Agent")
	tasks, errs := overviewExecutor(r)
	if errs != nil {
		c.ErrorHandler(w, r, errs)
	} else {
		c.SuccessHandler(w, r, tasks)
	}
}
