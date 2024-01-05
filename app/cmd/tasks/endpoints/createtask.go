package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/tasks/database"
	"logbook/config"
	"logbook/internal/web/reqs"
	"net/http"
)

type CreateTaskRequest struct {
	ParentOid database.ObjectiveId `json:"parent_oid"`
	ParentVid database.VersionId   `json:"parent_vid"`
	Text      string               `json:"text"`
}

func (ct CreateTaskRequest) validate() error {
	if !ct.ParentOid.Validate() {
		return fmt.Errorf("invalid value for 'parent' parameter")
	}
	return nil
}

type CreateTaskResponse struct {
	Created database.Objective   `json:"created"`
	Updated []database.Objective `json:"updated"`
}

// TODO: Check user input for script tags in order to prevent XSS attempts
func (e *Endpoints) CreateTask(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[CreateTaskRequest](r)
	if err != nil {
		log.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		log.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	o, err := e.app.CreateObjective(bq.ParentOid, bq.ParentVid)
	if err != nil {
		log.Println(fmt.Errorf("app.createObjective: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreateTaskResponse{
		Created: tasks[0],
		Updated: tasks[min(1, len(tasks)):],
	}
	if err := reqs.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (bq *CreateTaskRequest) Send() (*CreateTaskResponse, error) {
	return reqs.Send[CreateTaskRequest, CreateTaskResponse](config.ObjectivesCreate, bq)
}
