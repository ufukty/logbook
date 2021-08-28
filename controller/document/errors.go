package document

import (
	"encoding/json"
	"log"
	"net/http"

	"logbook/main/database"

	"github.com/google/uuid"
)

type ControllerError struct {
	httpStatusCode int
	helpMessage    string
}

var errorMapFromDatabaseToController map[error]ControllerError

var ErrNotSpecified ControllerError

var ErrNoResult ControllerError

var ErrInvalidInput ControllerError

var ErrEmptyUserId ControllerError

var ErrInvalidUserId ControllerError

var ErrEmptyDocumentId ControllerError

var ErrInvalidDocumentId ControllerError

var ErrEmptyTaskGroupId ControllerError

var ErrInvalidTaskGroupId ControllerError

var ErrEmptyDocumentName ControllerError

func initErrors() {

	ErrNotSpecified = ControllerError{
		httpStatusCode: http.StatusInternalServerError,
		helpMessage:    "Check back soon.",
	}
	ErrNoResult = ControllerError{
		httpStatusCode: http.StatusNotFound,
		helpMessage:    "No resources.",
	}
	ErrInvalidInput = ControllerError{
		httpStatusCode: http.StatusBadRequest,
		helpMessage:    "One or more paratemeters are invalid.",
	}
	ErrEmptyUserId = ControllerError{
		httpStatusCode: 404,
		helpMessage:    "",
	}
	ErrInvalidUserId = ControllerError{
		httpStatusCode: 404,
		helpMessage:    "",
	}
	ErrEmptyDocumentId = ControllerError{
		httpStatusCode: 404,
		helpMessage:    "",
	}
	ErrInvalidDocumentId = ControllerError{
		httpStatusCode: 404,
		helpMessage:    "",
	}
	ErrEmptyTaskGroupId = ControllerError{
		httpStatusCode: 404,
		helpMessage:    "",
	}
	ErrInvalidTaskGroupId = ControllerError{
		httpStatusCode: 404,
		helpMessage:    "",
	}
	ErrEmptyDocumentName = ControllerError{
		httpStatusCode: 404,
		helpMessage:    "A name for new document should be specified.",
	}

	errorMapFromDatabaseToController = map[error]ControllerError{
		database.ErrNotSpecified:       ErrNotSpecified,
		database.ErrNoResult:           ErrNoResult,
		database.ErrInvalidInput:       ErrInvalidInput,
		database.ErrEmptyDocumentId:    ErrEmptyDocumentId,
		database.ErrInvalidDocumentId:  ErrInvalidDocumentId,
		database.ErrEmptyTaskGroupId:   ErrEmptyTaskGroupId,
		database.ErrInvalidTaskGroupId: ErrInvalidTaskGroupId,
	}
}

func exportError(err error) ControllerError {
	return errorMapFromDatabaseToController[err]
}

func writeResponse(errorId string, err error, w http.ResponseWriter) {
	exportedError := exportError(err)
	json.NewEncoder(w).Encode(ResponseFields{
		Status:    exportedError.httpStatusCode,
		ErrorHint: exportedError.helpMessage,
		ErrorId:   errorId,
	})
	w.WriteHeader(exportedError.httpStatusCode)
}

func writeLog(
	errorId string,
	err error,
	endpoint string,
	w http.ResponseWriter,
	r *http.Request,
) {
	exportedError := exportError(err)
	byte_str, err := json.Marshal(LogFields{
		ErrorId:  errorId,
		Status:   exportedError.httpStatusCode,
		Endpoint: endpoint,
		Request:  r.Header,
	})
	if err != nil {
		log.Println("[WARNING] writeLog function can not print logs because of yaml.Marshall gives error.")
	}
	log.Println(string(byte_str))
}

func errorHandler(w http.ResponseWriter, r *http.Request, endpoint string, err error) {
	errorId := uuid.New().String()
	writeResponse(errorId, err, w)
	writeLog(errorId, err, endpoint, w, r)
}
