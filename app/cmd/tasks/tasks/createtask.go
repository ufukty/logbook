package tasks

import (
	"fmt"
	"log"
	"logbook/cmd/tasks/database"
	"logbook/internal/web/reqs"
	"net/http"
	"time"

	"github.com/jackc/pgtype"
)

type CreateTask struct {
	Parent database.Iid `json:"parent"`
	Text   string       `json:"text"`
}

func (ct CreateTask) Validate() error {
	if !ct.Parent.Validate() {
		return fmt.Errorf("invalid value for 'parent' parameter")
	}
	return nil
}

type CreateTaskR struct {
	Created database.Task   `json:"created"`
	Updated []database.Task `json:"updated"`
}

func (e *Endpoints) CreateTask(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[CreateTask](r)
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

	task := &database.Task{
		TaskId:        "",
		DocumentId:    "",
		ParentId:      "",
		Content:       "",
		Degree:        0,
		Depth:         0,
		CreatedAt:     time.Time{},
		CompletedAt:   pgtype.Date{},
		ReadyToPickUp: false,
	}

	tasks, err := e.db.CreateTask(task)
	if err != nil || len(tasks) == 0 {
		log.Println(fmt.Errorf("querying the db: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// check auth

	// creation of task
	// database.CreateTask(database.Task{
	// 	RevisionId:            "00000000-0000-0000-0000-000000000000",
	// 	OriginalCreatorUserId: "00000000-0000-0000-0000-000000000000",
	// 	ResponsibleUserId:     "00000000-0000-0000-0000-000000000000",
	// 	Content:               "Lorem ipsum dolor sit amet",
	// 	Notes:                 "Consectetur adipiscing elit",
	// })

	// creation of ownership role in PERM
	// database.CreatePermission(database.TaskPermission{
	// 	UserId: "00000000-0000-0000-0000-000000000000",
	// 	Role: "Role.Ownership",
	// })

	// check existence of super task

	// create link in TASK_LINK table

	// check permissions between task and user

	// create NewOperation

	// trigger task-props calculation

	bs := CreateTaskR{
		Created: tasks[0],
		Updated: tasks[min(1, len(tasks)):],
	}
	if err := reqs.WriteJsonResponse(bs, w); err != nil {
		log.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
