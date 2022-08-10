package responder

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

// Errors from [GET]/document/details
var (
	ErrDocumentIdInputIsNotValidUUID = errors.New("ErrDocumentIdInputIsNotValidUUID")
	ErrCreateTaskParseForm           = errors.New("ErrCreateTaskParseForm")
	ErrCreateTaskEmptyContent        = errors.New("ErrCreateTaskEmptyContent")
	ErrCreateTaskEmptyDocumentId     = errors.New("ErrCreateTaskEmptyDocumentId")
	ErrCreateTaskInvalidDocumentId   = errors.New("ErrCreateTaskInvalidDocumentId")
	ErrCreateTaskEmptyParentId       = errors.New("ErrCreateTaskEmptyParentId")
	ErrCreateTaskInvalidParentId     = errors.New("ErrCreateTaskInvalidParentId")
)

// Errors from [POST]/task
var (
	ErrUpdateParentParentCheck         = errors.New("UpdateParent faced with an error when running db.GetTaskByTaskId(task.ParentId) to check if parent task id is valid")
	ErrUpdateParentSaveChanges         = errors.New("UpdateParent faced with an error when trying to save changes into database")
	ErrUpdateParentMaximumDepthReached = errors.New("ErrUpdateParentMaximumDepthReached")
	ErrUpdateParentNextTaskCheck       = errors.New("UpdateParent faced with an error when trying to check next child task")
	ErrTaskCreateUpdateParents         = errors.New("Task/createExecutor faced with an error while trying to ")
)

// Errors from [GET]/document/overview/chronological
var (
	ErrOffsetInputIsNotValidInteger   = errors.New("ErrOffsetInputIsNotValidInteger")
	ErrOffsetInputIsNotInAllowedRange = errors.New("ErrOffsetInputIsNotInAllowedRange")
	ErrLimitInputIsNotValidInteger    = errors.New("ErrLimitInputIsNotValidInteger")
	ErrLimitInputIsNotInAllowedRange  = errors.New("ErrLimitInputIsNotInAllowedRange")
)

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
	ErrorHint     string      `json:"error_hint" yaml:"error_hint"`
	ErrorStack    error       `json:"error_stack" yaml:"error_stack"`
	RequestHeader interface{} `json:"request_header" yaml:"request_header"`
	RequestForm   interface{} `json:"request_form" yaml:"request_form"`
	Endpoint      string      `json:"endpoint" yaml:"endpoint"`
}

func InternalErrorHandler(
	r *http.Request,
	incidentId string,
	statusCode int,
	errorMessageForResponse string,
	errStackForLogs error,
	endpoint string,
) {
	byte_str, err := yaml.Marshal(ControllerLoggingFields{
		IncidentId:    incidentId,
		ErrorHint:     errorMessageForResponse,
		ErrorStack:    errStackForLogs,
		RequestHeader: r.Header,
		RequestForm:   r.PostForm,
		Status:        statusCode,
		Endpoint:      endpoint,
	})
	if err != nil {
		log.Println("[WARNING] InternalErrorHandler function can not print logs because of yaml.Marshal(ControllerLoggingFields{...}) gives error.")
	}
	log.Println(string(byte_str))
}

func PublicFacingErrorHandler(
	w http.ResponseWriter,
	incidentId string,
	statusCode int,
	errorMessageForResponse string,
) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ControllerResponseFields{
		Status:     statusCode,
		ErrorHint:  errorMessageForResponse,
		IncidentId: incidentId,
	})
}

func ErrorHandler(
	w http.ResponseWriter,
	r *http.Request,
	statusCode int,
	errorMessageForResponse string,
	errStackForLogs error,
) {
	incidentId := uuid.New().String()
	PublicFacingErrorHandler(w, incidentId, statusCode, errorMessageForResponse)

	var endpoint string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		endpoint = details.Name()
	} else {
		endpoint = "could not traced the endpoint"
	}

	InternalErrorHandler(r, incidentId, statusCode, errorMessageForResponse, errStackForLogs, endpoint)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
