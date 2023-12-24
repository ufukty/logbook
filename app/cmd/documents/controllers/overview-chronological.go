package document

import (
	// db "logbook/cmd/tasks/database"
	// c "logbook/cmd/tasks/controller"
	// e "logbook/internal/errors"
	// v "logbook/internal/validate"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	OFFSET_ALLOWED_RANGE_MIN = 0
	OFFSET_ALLOWED_RANGE_MAX = 10000000
	LIMIT_ALLOWED_RANGE_MIN  = 1
	LIMIT_ALLOWED_RANGE_MAX  = 100000
)

type CDocumentOverviewChronological struct {
	userInput struct {
		DocumentId string
		Limit      int
		Offset     int
	}
}

func (s *CDocumentOverviewChronological) sanitizer(r *http.Request) []error {
	vars := mux.Vars(r)

	documentId := vars["document_id"]
	if !v.IsValidUUID(documentId) {
		return []error{c.ErrDocumentIdInputIsNotValidUUID}
	}

	offset, err := strconv.Atoi(vars["offset"])
	if err != nil {
		return []error{err, c.ErrOffsetInputIsNotValidInteger}
	}
	if offset < OFFSET_ALLOWED_RANGE_MIN || OFFSET_ALLOWED_RANGE_MAX < offset {
		return []error{c.ErrOffsetInputIsNotInAllowedRange}
	}

	limit, err := strconv.Atoi(vars["limit"])
	if err != nil {
		return []error{err, c.ErrLimitInputIsNotValidInteger}
	}
	if limit < LIMIT_ALLOWED_RANGE_MIN || LIMIT_ALLOWED_RANGE_MAX < limit {
		return []error{c.ErrLimitInputIsNotInAllowedRange}
	}

	s.userInput.DocumentId = documentId
	s.userInput.Limit = limit
	s.userInput.Offset = offset

	return nil
}

func (s *CDocumentOverviewChronological) executor(r *http.Request) ([]db.Task, *e.Error) {
	var (
		tasks []db.Task
		errs  []error
	)

	errs = s.sanitizer(r)
	if errs != nil {
		return nil, e.New("Check your inputs.", http.StatusBadRequest, errs)
	}

	// create document table record
	tasks, errs = db.GetChronologicalViewItems(s.userInput.DocumentId, s.userInput.Limit, s.userInput.Offset)
	if errs != nil {
		return nil, e.New(errs)
	}

	return tasks, nil
}

func (s *CDocumentOverviewChronological) Responder(w http.ResponseWriter, r *http.Request) {
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
