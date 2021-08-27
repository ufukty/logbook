package document

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type ControllerError struct {
	httpStatusCode int
	helpMessage    string
}

// var errorMapFromDatabaseToController map[error]ControllerError

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

	// 	errorMapFromDatabaseToController = map[error]ControllerError{
	// 		database.ErrNotSpecified:       ErrNotSpecified,
	// 		database.ErrNoResult:           ErrNoResult,
	// 		database.ErrInvalidInput:       ErrInvalidInput,
	// 		database.ErrEmptyUserId:        ErrEmptyUserId,
	// 		database.ErrInvalidUserId:      ErrInvalidUserId,
	// 		database.ErrEmptyDocumentId:    ErrEmptyDocumentId,
	// 		database.ErrInvalidDocumentId:  ErrInvalidDocumentId,
	// 		database.ErrEmptyTaskGroupId:   ErrEmptyTaskGroupId,
	// 		database.ErrInvalidTaskGroupId: ErrInvalidTaskGroupId,
	// 	}
}

// func exportError(err error) ControllerError {
// 	return errorMapFromDatabaseToController[err]
// }

func writeResponse(eventId string, err ControllerError, requestParameters map[string]string, w http.ResponseWriter) {
	w.WriteHeader(err.httpStatusCode)
	json.NewEncoder(w).Encode(struct {
		EventId           string            `json:"event_id"`
		Status            bool              `json:"status"`
		Hint              string            `json:"hint"`
		RequestParameters map[string]string `json:"request_parameters"`
	}{
		EventId:           eventId,
		Status:            false,
		Hint:              err.helpMessage,
		RequestParameters: requestParameters,
	})
}

func writeLog(eventId string, controllerError ControllerError, originalError error, requestParameters map[string]string, w http.ResponseWriter) {
	byte_str, err := json.Marshal(struct {
		EventId           string            `yaml:"event_id"`
		ControllerError   ControllerError   `yaml:"controller_error"`
		OriginalError     error             `yaml:"original_error"`
		RequestParameters map[string]string `yaml:"request_parameters"`
	}{
		EventId:           eventId,
		ControllerError:   controllerError,
		OriginalError:     originalError,
		RequestParameters: requestParameters,
	})
	if err != nil {
		log.Println("[WARNING] writeLog function can not print logs because of yaml.Marshall gives error.")
	}
	log.Println(string(byte_str))
}

func errorHandler(
	controllerError ControllerError,
	originalError error,
	requestParameters map[string]string,
	w http.ResponseWriter,
) {
	eventId := uuid.New().String()
	writeResponse(eventId, controllerError, requestParameters, w)
	writeLog(eventId, controllerError, originalError, requestParameters, w)
}
