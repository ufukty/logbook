package document

import (
	c "logbook/controller"
	e "logbook/controller/utilities/errors"
	db "logbook/database"
	"net/http"

	"github.com/gorilla/mux"
)

func sanitizeUserInput(r *http.Request) (*db.Document, *e.Error) {

	uriVars := mux.Vars(r)

	documentId := uriVars["document_id"]

	if documentId == "" {
		return nil, e.New(`Check your input for 'document_id'.`, http.StatusBadRequest)
	}

	document, errs := db.GetDocumentByDocumentId(documentId)
	if errs != nil {
		return nil, e.New(errs)
	}

	return &document, nil
}

func detailsExecutor(r *http.Request) (*db.Document, *e.Error) {
	document, err := sanitizeUserInput(r)
	if err != nil {
		return nil, err
	}

	taskGroups, errs := db.GetTaskGroupsByDocumentId(document.DocumentId)
	if errs != nil {
		return nil, e.New(errs)
	}
	document.TaskGroups = taskGroups

	for _, taskGroup := range taskGroups {
		switch taskGroup.TaskGroupType {
		case db.Archive:
			passthrough
		}
	}

	return document, nil
}

func Details(w http.ResponseWriter, r *http.Request) {
	// ipAddress := (*r).RemoteAddr
	// userAgent := (*r).Header.Get("User-Agent")
	document, errs := detailsExecutor(r)
	if errs != nil {
		c.ErrorHandler(w, r, errs)
	} else {
		c.SuccessHandler(w, r, document)
	}
}
