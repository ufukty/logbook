package tasks

import (
	"fmt"
	"log"
	"logbook/cmd/tasks/database"
	"logbook/internal/web/reqs"
	"net/http"
)

type ListTasks struct {
	Root database.Iid `url:"root"`
}

type ListTasksR struct {
	Tasks []database.Task `json:"tasks"`
}

func (rq ListTasks) Validate() error {
	if !rq.Root.Validate() {
		return fmt.Errorf("invalid iid for root")
	}
	return nil
}

func (e Endpoints) ListTasks(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[ListTasks](r)
	if err != nil {
		log.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := bq.Validate(); err != nil {
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

	bs := ListTasksR{
		Tasks: tasks,
	}
	if err := reqs.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
