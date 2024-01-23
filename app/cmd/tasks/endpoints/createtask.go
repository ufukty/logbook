package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/tasks/app"
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
	Update []database.Ovid `json:"update"`
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

	act := app.CreateObjectiveAction{
		ParentOid: bq.ParentOid,
		ParentVid: bq.ParentVid,
		Content:   bq.Text,
		Creator:   "00000000-0000-0000-0000-000000000000", // TODO: check auth header
	}

	o, err := e.app.CreateObjective(act)
	if err != nil {
		log.Println(fmt.Errorf("app.createObjective: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreateTaskResponse{
		Created: o.Oid,
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
