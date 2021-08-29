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

var (
	ErrTaskGroupIdCheck             = errors.New("ErrTaskGroupIdCheck")
	ErrCreateTaskCall               = errors.New("ErrCreateTaskCall")
	ErrSterilizeUserInputDegreeInt  = errors.New("ErrSterilizeUserInputDegreeInt")
	ErrSterilizeUserInputDepthInt   = errors.New("ErrSterilizeUserInputDepthInt")
	ErrSterilizeUserInputTaskStatus = errors.New("ErrSterilizeUserInputTaskStatus")
	ErrSterilizeUserInputParseForm  = errors.New("ErrSterilizeUserInputParseForm")
)

var httpErrorMapping = map[error]int{
	ErrTaskGroupIdCheck:             http.StatusNotAcceptable,
	ErrCreateTaskCall:               http.StatusNotAcceptable,
	ErrSterilizeUserInputDegreeInt:  http.StatusNotAcceptable,
	ErrSterilizeUserInputDepthInt:   http.StatusNotAcceptable,
	ErrSterilizeUserInputTaskStatus: http.StatusNotAcceptable,
	ErrSterilizeUserInputParseForm:  http.StatusNotAcceptable,
}

var httpHintMapping = map[error]string{
	ErrTaskGroupIdCheck:             "Task Group Id is invalid. Check it again.",
	ErrCreateTaskCall:               "There was an error when attempted to create a task.",
	ErrSterilizeUserInputDegreeInt:  "Check your inputs or try again later.",
	ErrSterilizeUserInputDepthInt:   "Check your inputs or try again later.",
	ErrSterilizeUserInputTaskStatus: "Check your inputs or try again later.",
	ErrSterilizeUserInputParseForm:  "Check your inputs or try again later.",
}

// Used for both error and success messages
// But only for rendering public http response
type ControllerResponseFields struct {
	Status     int         `json:"status" yaml:"status"`
	IncidentId string      `json:"incident_id" yaml:"incident_id"`
	ErrorHint  string      `json:"error_hint" yaml:"error_hint"`
	Resource   interface{} `json:"resource" yaml:"resource"`
}

// Only used by ErrorHandler to hold error information
type ControllerError struct {
	Wrapper    error `json:"wrapper" yaml:"wrapper"`
	Underlying error `json:"underlying" yaml:"underlying"`
}

// Only used by ErrorHandler to hold error information
type ControllerErrorStrings struct {
	Wrapper    string `json:"wrapper" yaml:"wrapper"`
	Underlying string `json:"underlying" yaml:"underlying"`
}

// Used for both error and success messages
// But only for writing internal logs
type ControllerLoggingFields struct {
	Status        int                    `json:"status" yaml:"status"`
	IncidentId    string                 `json:"incident_id" yaml:"incident_id"`
	Error         ControllerErrorStrings `json:"controller_error" yaml:"controller_error"`
	RequestHeader interface{}            `json:"request_header" yaml:"request_header"`
	RequestForm   interface{}            `json:"request_form" yaml:"request_form"`
	Endpoint      string                 `json:"endpoint" yaml:"endpoint"`
}

func serializeControllerError(ce ControllerError) ControllerErrorStrings {
	return ControllerErrorStrings{
		Wrapper:    ce.Wrapper.Error(),
		Underlying: ce.Underlying.Error(),
	}
}

func InternalErrorHandler(
	w http.ResponseWriter,
	r *http.Request,
	incidentId string,
	controllerError ControllerError,
	endpoint string,
) {
	byte_str, err := yaml.Marshal(ControllerLoggingFields{
		IncidentId:    incidentId,
		Error:         serializeControllerError(controllerError),
		RequestHeader: r.Header,
		RequestForm:   r.PostForm,
		Status:        httpErrorMapping[controllerError.Wrapper],
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
	controllerError ControllerError,
) {
	w.WriteHeader(httpErrorMapping[controllerError.Wrapper])
	json.NewEncoder(w).Encode(ControllerResponseFields{
		Status:     httpErrorMapping[controllerError.Wrapper],
		ErrorHint:  httpHintMapping[controllerError.Wrapper],
		IncidentId: incidentId,
	})
}

func ErrorHandler(
	w http.ResponseWriter,
	r *http.Request,
	controllerError ControllerError,
) {
	errorId := uuid.New().String()
	PublicFacingErrorHandler(w, errorId, controllerError)
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		InternalErrorHandler(w, r, errorId, controllerError, details.Name())
	}
	panic("ErrorHandler called panic.")
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
