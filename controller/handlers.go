package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

// Errors from [POST]/document
var (
	ErrCreateDocumentCreateCreationPhase = errors.New("CreateDocument faced with an error when tried to create document with task groups")
	ErrCreateDocumentSelectionPhase      = errors.New("CreateDocument faced with an error when tried to access created task groups")
)

var errorResponseCodes = map[error]int{
	ErrCreateDocumentCreateCreationPhase: http.StatusNoContent,
	ErrCreateDocumentSelectionPhase:      http.StatusNoContent,
}

var errorHints = map[error]string{
	ErrCreateDocumentCreateCreationPhase: "Try again later.",
	ErrCreateDocumentSelectionPhase:      "Try again later.",
}

// Errors from [POST]/task
var (
	ErrTaskGroupIdCheck             = errors.New("ErrTaskGroupIdCheck")
	ErrCreateTaskCall               = errors.New("ErrCreateTaskCall")
	ErrSterilizeUserInputDegreeInt  = errors.New("ErrSterilizeUserInputDegreeInt")
	ErrSterilizeUserInputDepthInt   = errors.New("ErrSterilizeUserInputDepthInt")
	ErrSterilizeUserInputTaskStatus = errors.New("ErrSterilizeUserInputTaskStatus")
	ErrSterilizeUserInputParseForm  = errors.New("ErrSterilizeUserInputParseForm")
	ErrParentCheck                  = errors.New("ErrParentCheck")
	ErrCreateTaskUpdateParent       = errors.New("ErrCreateTaskUpdateParent")
	ErrMaximumDepthReached          = errors.New("ErrMaximumDepthReached")
	ErrNextTaskCheck                = errors.New("ErrNextTaskCheck")
)

var httpErrorMapping = map[error]int{
	ErrTaskGroupIdCheck:             http.StatusNotAcceptable,
	ErrCreateTaskCall:               http.StatusNotAcceptable,
	ErrSterilizeUserInputDegreeInt:  http.StatusNotAcceptable,
	ErrSterilizeUserInputDepthInt:   http.StatusNotAcceptable,
	ErrSterilizeUserInputTaskStatus: http.StatusNotAcceptable,
	ErrSterilizeUserInputParseForm:  http.StatusNotAcceptable,
	ErrParentCheck:                  http.StatusNotAcceptable,
	ErrCreateTaskUpdateParent:       http.StatusNotAcceptable,
	ErrMaximumDepthReached:          http.StatusNotAcceptable,
	ErrNextTaskCheck:                http.StatusNotAcceptable,
}

var httpHintMapping = map[error]string{
	ErrTaskGroupIdCheck:             "task-group-id is invalid. Check it again.",
	ErrCreateTaskCall:               "There was an error when attempted to create a task.",
	ErrSterilizeUserInputDegreeInt:  "'Degree' is invalid. Check your inputs or try again later.",
	ErrSterilizeUserInputDepthInt:   "'depth' is invalid. Check your inputs or try again later.",
	ErrSterilizeUserInputTaskStatus: "'task-status' is invalid. Check your inputs or try again later.",
	ErrSterilizeUserInputParseForm:  "HTTP 'form' oject is invalid. Check your inputs or try again later.",
	ErrParentCheck:                  "No such task as in parent-id.",
	ErrCreateTaskUpdateParent:       "Relation with parent task may not updated properly.",
	ErrMaximumDepthReached:          "Maximum depth is reached.",
	ErrNextTaskCheck:                "Document could not processed.",
}

// Used for both error and success messages
// But only for rendering public http response
type ControllerResponseFields struct {
	Status     int         `json:"status" yaml:"status"`
	IncidentId string      `json:"incident_id" yaml:"incident_id"`
	ErrorHint  string      `json:"error_hint" yaml:"error_hint"`
	Resource   interface{} `json:"resource" yaml:"resource"`
}

// Used for both error and success messages
// But only for writing internal logs
type ControllerLoggingFields struct {
	Status        int         `json:"status" yaml:"status"`
	IncidentId    string      `json:"incident_id" yaml:"incident_id"`
	ErrorStack    []string    `json:"error_stack" yaml:"error_stack"`
	RequestHeader interface{} `json:"request_header" yaml:"request_header"`
	RequestForm   interface{} `json:"request_form" yaml:"request_form"`
	Endpoint      string      `json:"endpoint" yaml:"endpoint"`
}

func serializeControllerError(errs []error) []string {
	errs_str := []string{}
	for _, err := range errs {
		errs_str = append(errs_str, err.Error())
	}
	return errs_str
}

func InternalErrorHandler(
	r *http.Request,
	incidentId string,
	errs []error,
	endpoint string,
	responseCode int,
	errorHint string,
) {
	byte_str, err := yaml.Marshal(ControllerLoggingFields{
		IncidentId:    incidentId,
		ErrorStack:    serializeControllerError(errs),
		RequestHeader: r.Header,
		RequestForm:   r.PostForm,
		Status:        responseCode,
		Endpoint:      endpoint,
	})
	if err != nil {
		log.Println("[WARNING] writeLog function can not print logs because of yaml.Marshall gives error.")
	}
	log.Println(string(byte_str))
}

func PublicFacingErrorHandler(
	w http.ResponseWriter,
	incidentId string,
	errs []error,
	responseCode int,
	errorHint string,
) {
	w.WriteHeader(responseCode)
	json.NewEncoder(w).Encode(ControllerResponseFields{
		Status:     responseCode,
		ErrorHint:  errorHint,
		IncidentId: incidentId,
	})
}

func ErrorHandler(
	w http.ResponseWriter,
	r *http.Request,
	errs []error,
) {
	errorId := uuid.New().String()

	var responseCode int
	if code, ok := errorResponseCodes[errs[len(errs)-1]]; ok {
		responseCode = code
	} else {
		responseCode = http.StatusInternalServerError
	}

	var errorHint string
	if code, ok := errorHints[errs[len(errs)-1]]; ok {
		errorHint = code
	} else {
		errorHint = "Unexpected error. Try again later."
	}

	PublicFacingErrorHandler(w, errorId, errs, responseCode, errorHint)

	var endpoint string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		endpoint = details.Name()
	} else {
		endpoint = "could not traced the endpoint"
	}

	InternalErrorHandler(r, errorId, errs, endpoint, responseCode, errorHint)
}

func ResponseHandler(
	w http.ResponseWriter,
	r *http.Request,
	controllerResponse ControllerResponseFields,
) {
	json.NewEncoder(w).Encode(controllerResponse)
}

func InternalSuccessHandler(endpoint string) {
	log.Println("Request processed @", endpoint)
}

func SuccessHandler(
	w http.ResponseWriter,
	r *http.Request,
	resource interface{},
) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		InternalSuccessHandler(details.Name())
	}
	ResponseHandler(w, r, ControllerResponseFields{
		Status:   http.StatusOK,
		Resource: resource,
	})
}
