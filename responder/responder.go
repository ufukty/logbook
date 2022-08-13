package responder

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
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
	ErrorStack    string      `json:"error_stack" yaml:"error_stack"`
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
		Status:        statusCode,
		IncidentId:    incidentId,
		ErrorHint:     errorMessageForResponse,
		ErrorStack:    fmt.Sprint(errStackForLogs),
		RequestHeader: r.Header,
		RequestForm:   r.PostForm,
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
