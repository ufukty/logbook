package app

import (
	"fmt"
	"log"
	"logbook/cmd/tags/endpoints"
	"logbook/cmd/tasks/database"
	"net/http"

	"github.com/google/uuid"
)

func (a *App) ApplyActionsOnObjective(o *database.Objective, as []database.Action) error {
	// apply all actions in one version number change
	for _, a := range as {
		switch a.Summary {
		// case db.TASK_CREATE:

		// case db.TASK_REORDER:

		// case db.TASK_DELETE:

		case database.ObjectiveTextAssign:

		case database.ObjectiveMarkComplete:

		case database.ObjectiveUnmarkComplete:

		case database.ObjectiveNoteAssign:

		case database.CollaborationAssign:

		case database.CollaborationUnassign:

		case database.COLLABORATION_RESTRICT:

		case database.COLLABORATION_DERESTRICT:

		case database.COLLABORATION_CHANGE_ROLE:

		case database.HISTORY_ROLLBACK:

		case database.HISTORY_FASTFORWARD:

		}
	}
	return nil
}

// proposals can have multiple actions
func (a *App) UpdateObjective(oid database.ObjectiveId, vid database.VersionId, as []database.Action) error {
	nextVid := uuid.New()

	o, err := a.db.SelectObjective(oid, vid)
	if err != nil {
		return fmt.Errorf("getting objective for oid: %w", err)
	}

	if err := a.ApplyActionsOnObjective(o.Clone(), as); err != nil {
		return fmt.Errorf("applying action list on objective: %w", err)
	}

	createNextVersionOfParent := func(oid database.ObjectiveId) database.Objective {
		links := a.db.SuperLinksVersioned(oid, vid)
		for _, link := range links {
			if link.Type == nil {
			}
		}
		// TODO: same version of update sibling
	}
	updateChildren := func(oid database.ObjectiveId, vid database.VersionId) {

	}

	bq := endpoints.TagAssignRequest{}
	bs, err := bq.Send()
	if err != nil {
		return fmt.Errorf("copying tag records to new version of the task: %w", err)
	}

	return nil
}

func (a *App) ComposeView(root database.ObjectiveId, vid database.VersionId) (any, error) {

}

// TODO: Create new version of parent (and whole ancestry)
// TODO: Turn the parent objective into a goal if it is currently a task
func (a *App) CreateObjective() {
	// TODO: get version number to create updated versions of ancestry with it

	a.db.SelectObjective()

	// TODO:
	o := database.Objective{
		Vid:     "",
		Based:   "",
		Type:    "",
		Content: "",
		Creator: "",
	}

	// TODO: link to parent

	o, err = a.db.InsertObjective(o)
	if err != nil {
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
}
