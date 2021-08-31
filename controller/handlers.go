package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

	e "logbook/main/controller/utilities/errors"
)

// Errors from [POST]/document
var (
	ErrCreateDocumentCreateCreationPhase = errors.New("CreateDocument faced with an error when tried to create document with task groups")
	ErrCreateDocumentSelectionPhase      = errors.New("CreateDocument faced with an error when tried to access created task groups")
)

// Errors from [POST]/task
var (
	ErrTaskGroupIdCheck                      = errors.New("ErrTaskGroupIdCheck")
	ErrTaskCreateCreateTaskCall              = errors.New("ErrCreateTaskCall")
	ErrTaskCreateSanitizeDegreeConvertion    = errors.New("problem with converting 'degree' to integer")
	ErrTaskCreateSanitizeDegreeNegativeValue = errors.New("problem with converting 'degree' to integer")
	ErrTaskCreateSanitizeDepthConvertion     = errors.New("problem with converting 'depth' to integer")
	ErrTaskCreateSanitizeDepthNegativeValue  = errors.New("problem with converting 'depth' to integer")
	ErrTaskCreateSanitizeTaskStatus          = errors.New("problem with converting 'task_status' to TaskStatus type")
	ErrTaskCreateSanitizeParseForm           = errors.New("problem with parsing http.request via ParseForm()")
	ErrTaskCreateSanitizeParentId            = errors.New("problem with validating 'parent_id'")
	ErrTaskCreateSanitizeTaskGroupId         = errors.New("problem with validating 'task_group_id'")
	ErrUpdateParentParentCheck               = errors.New("UpdateParent faced with an error when running db.GetTaskByTaskId(task.ParentId) to check if parent task id is valid")
	ErrUpdateParentSaveChanges               = errors.New("UpdateParent faced with an error when trying to save changes into database")
	ErrUpdateParentMaximumDepthReached       = errors.New("ErrUpdateParentMaximumDepthReached")
	ErrUpdateParentNextTaskCheck             = errors.New("UpdateParent faced with an error when trying to check next child task")
	ErrTaskCreateSanitize                    = errors.New("Task/createExecutor faced with an error while trying to sanitize user input")
	ErrTaskCreateTaskGroupIdCheck            = errors.New("Task/createExecutor faced with an error while trying to ")
	ErrTaskCreateParentCheck                 = errors.New("Task/createExecutor faced with an error while trying to ")
	ErrTaskCreateUpdateParents               = errors.New("Task/createExecutor faced with an error while trying to ")
)

var errorResponseCodes = map[error]int{
	ErrCreateDocumentCreateCreationPhase: http.StatusNoContent,
	ErrCreateDocumentSelectionPhase:      http.StatusNoContent,
}

var errorHints = map[error]string{
	ErrCreateDocumentCreateCreationPhase: "Try again later.",
	ErrCreateDocumentSelectionPhase:      "Try again later.",
	ErrUpdateParentMaximumDepthReached:   "You can't create new dependency to this task due to maximum depth limit.",
}

// var httpErrorMapping = map[error]int{
// 	ErrTaskGroupIdCheck:                http.StatusNotAcceptable,
// 	ErrTaskCreateCreateTaskCall:        http.StatusNotAcceptable,
// 	ErrSterilizeUserInputDegreeInt:     http.StatusNotAcceptable,
// 	ErrSterilizeUserInputDepthInt:      http.StatusNotAcceptable,
// 	ErrSterilizeUserInputTaskStatus:    http.StatusNotAcceptable,
// 	ErrSterilizeUserInputParseForm:     http.StatusNotAcceptable,
// 	ErrUpdateParentParentCheck:         http.StatusNotAcceptable,
// 	ErrUpdateParentSaveChanges:         http.StatusNotAcceptable,
// 	ErrUpdateParentMaximumDepthReached: http.StatusNotAcceptable,
// 	ErrUpdateParentNextTaskCheck:       http.StatusNotAcceptable,
// }

// var httpHintMapping = map[error]string{
// 	ErrTaskGroupIdCheck:             "task-group-id is invalid. Check it again.",
// 	ErrTaskCreateCreateTaskCall:     "There was an error when attempted to create a task.",
// 	ErrSterilizeUserInputDegreeInt:  "'Degree' is invalid. Check your inputs or try again later.",
// 	ErrSterilizeUserInputDepthInt:   "'depth' is invalid. Check your inputs or try again later.",
// 	ErrSterilizeUserInputTaskStatus: "'task-status' is invalid. Check your inputs or try again later.",
// 	ErrSterilizeUserInputParseForm:  "HTTP 'form' oject is invalid. Check your inputs or try again later.",
// 	// ErrParentCheck:                  "No such task as in parent-id.",
// 	// ErrCreateTaskUpdateParent:       "Relation with parent task may not updated properly.",
// 	// ErrMaximumDepthReached:          "Maximum depth is reached.",
// 	// ErrNextTaskCheck:                "Document could not processed.",
// }

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
	errs *e.Error,
	endpoint string,
) {
	byte_str, err := yaml.Marshal(ControllerLoggingFields{
		IncidentId:    incidentId,
		ErrorHint:     errs.HttpResponseHint,
		ErrorStack:    serializeControllerError(errs.ErrorTrace),
		RequestHeader: r.Header,
		RequestForm:   r.PostForm,
		Status:        errs.HttpResponseCode,
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
	errs *e.Error,
) {
	w.WriteHeader(errs.HttpResponseCode)
	json.NewEncoder(w).Encode(ControllerResponseFields{
		Status:     errs.HttpResponseCode,
		ErrorHint:  errs.HttpResponseHint,
		IncidentId: incidentId,
	})
}

func ErrorHandler(
	w http.ResponseWriter,
	r *http.Request,
	errs *e.Error,
) {
	errorId := uuid.New().String()
	PublicFacingErrorHandler(w, errorId, errs)

	var endpoint string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		endpoint = details.Name()
	} else {
		endpoint = "could not traced the endpoint"
	}

	InternalErrorHandler(r, errorId, errs, endpoint)
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
