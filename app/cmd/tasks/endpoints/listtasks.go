package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/tasks/database"
	"logbook/internal/web/reqs"
	"net/http"
)

type ListObjectivesRequest struct {
	Root database.ObjectiveId `url:"root"`
}

type ListObjectivesResponse struct {
	Tasks []database.Objective `json:"tasks"`
}

func (rq ListObjectivesRequest) validate() error {
	if !rq.Root.Validate() {
		return fmt.Errorf("invalid iid for root")
	}
	return nil
}

func (e Endpoints) ListObjectives(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[ListObjectivesRequest](r)
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

	tasks, err := e.db.ListTasksInDocument(string(bq.Root))
	if err != nil {
		log.Println(fmt.Errorf("querying db: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := ListObjectivesResponse{
		Tasks: tasks,
	}
	if err := reqs.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
