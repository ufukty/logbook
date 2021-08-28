package document

import (
	"encoding/json"
	"log"
	"logbook/main/database"
	"net/http"

	"gopkg.in/yaml.v3"
)

type ResponseFields struct {
	Status    int         `json:"status" yaml:"status"`
	Resource  interface{} `json:"request" yaml:"request"`
	ErrorHint string      `json:"error_hint" yaml:"error_hint"`
	ErrorId   string      `json:"error_id" yaml:"error_id"`
}

type LogFields struct {
	Status   int         `json:"status" yaml:"status"`
	Endpoint string      `json:"endpoint" yaml:"endpoint"`
	Request  interface{} `json:"request" yaml:"request"`
	ErrorId  string      `json:"error_id" yaml:"error_id"`
}

func responseHandler(
	w http.ResponseWriter,
	r *http.Request,
	endpoint string,
	requestedResource interface{},
) {
	err := json.NewEncoder(w).Encode(ResponseFields{
		Status:   http.StatusOK,
		Resource: requestedResource,
	})
	if err != nil {
		log.Println("[WARNING] responseHandler function can not print HTTP responses because of json.NewEncoder().Encode() gives error.")
	}
	bytes, err := yaml.Marshal(LogFields{
		Status:   http.StatusOK,
		Endpoint: endpoint,
		Request:  r.Header,
	})
	if err != nil {
		log.Println("[WARNING] responseHandler function can not print logs because of yaml.Marshall gives error.")
	}
	log.Println(string(bytes))
}

func Create(w http.ResponseWriter, r *http.Request) {
	// ipAddress := (*r).RemoteAddr
	// userAgent := (*r).Header.Get("User-Agent")

	var (
		document database.Document
		err      error
	)

	// create document table record
	document, err = database.CreateDocumentWithTaskGroups(database.Document{})
	if err != nil {
		errorHandler(w, r, "Document/Create", err)
	}

	taskGroups, err := database.GetTaskGroupsByDocumentId(document.DocumentId)
	if err != nil {
		errorHandler(w, r, "Document/Create", err)
	}

	document.TaskGroups = taskGroups

	responseHandler(w, r, "Document/Create", document)
}
