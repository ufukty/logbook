package document

import (
	c "logbook/main/controller"
	e "logbook/main/controller/utilities/errors"
	v "logbook/main/controller/utilities/validate"
	db "logbook/main/database"
	"net/http"

	m "logbook/main/models"

	"github.com/gorilla/mux"
)

type CDocumentOverviewHierarchical struct {
	userInput struct {
		UserId     m.UUID
		DocumentId m.UUID
		Limit      int // TODO:
		Offset     int // TODO:
	}
}

func (s *CDocumentOverviewHierarchical) sanitizer(r *http.Request) (string, []error) {
	vars := mux.Vars(r)
	userId := vars["user_id"] // FIXME:
	documentId := vars["document_id"]
	if !v.IsValidUUID(documentId) {
		return "", []error{c.ErrDocumentIdInputIsNotValidUUID}
	}
	if !v.IsValidUUID(documentId) {
		return "", []error{c.ErrDocumentIdInputIsNotValidUUID}
	}
	return documentId, nil
}

func (s *CDocumentOverviewHierarchical) executor(r *http.Request) ([]db.Task, *e.Error) {
	var (
		tasks []db.Task
		errs  []error
	)

	documentId, errs := s.sanitizer(r)
	if errs != nil {
		return nil, e.New("Check your inputs.", http.StatusBadRequest, errs)
	}

	// create document table record
	requestedDocument := db.RequestedDocument{
		UserId:     userId,
		DocumentId: m.UUID(documentId),
	}
	tasks, errs = requestedDocument.GetDetails()
	if errs != nil {
		return nil, e.New(errs)
	}

	return tasks, nil
}

func (s *CDocumentOverviewHierarchical) Responder(w http.ResponseWriter, r *http.Request) {
	// ipAddress := (*r).RemoteAddr
	// userAgent := (*r).Header.Get("User-Agent")
	tasks, errs := s.executor(r)
	if errs != nil {
		c.ErrorHandler(w, r, errs)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		c.SuccessHandler(w, r, tasks)
	}
}
