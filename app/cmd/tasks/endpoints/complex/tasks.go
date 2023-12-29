package complex

import (
	"fmt"
	db "logbook/cmd/tasks/database"

	"github.com/google/uuid"
)

func ApplyActionsOnObjective(o *db.Objective, as []db.Action) error {
	// apply all actions in one version number change
	for _, a := range as {
		switch a.Summary {
		// case db.TASK_CREATE:

		// case db.TASK_REORDER:

		// case db.TASK_DELETE:

		case db.ObjectiveTextAssign:

		case db.ObjectiveMarkComplete:

		case db.ObjectiveUnmarkComplete:

		case db.ObjectiveNoteCreate:

		case db.ObjectiveNoteAssign:

		case db.ObjectiveNoteDelete:

		case db.CollaborationAssign:

		case db.CollaborationUnassign:

		case db.COLLABORATION_RESTRICT:

		case db.COLLABORATION_DERESTRICT:

		case db.COLLABORATION_CHANGE_ROLE:

		case db.HISTORY_ROLLBACK:

		case db.HISTORY_FASTFORWARD:

		}
	}
	return nil
}

// proposals can have multiple actions
func UpdateObjective(db *db.Database, oid db.ObjectiveId, vid db.VersionId, as []db.Action) error {
	nextVid := uuid.New()

	o, err := db.GetObjective(oid)
	if err != nil {
		return fmt.Errorf("getting objective for oid: %w", err)
	}

	if err := ApplyActionsOnObjective(o.Clone(), as); err != nil {
		return fmt.Errorf("applying action list on objective: %w", err)
	}

	createNextVersionOfParent := func(oid db.ObjectiveId) db.Objective {
		links := SuperLinks(oid, vid)
		for _, link := range links {
			if link.Type == nil {
			}
		}
		// TODO: same version of update sibling
	}
	updateChildren := func(oid db.ObjectiveId, vid db.VersionId) {

	}

	return nil
}

func GetLastVersionInBranch(vid db.VersionId) (db.VersionId, error) {
	prevVid, err := db.GetPreviousVersion(vid)
	if err != nil {
		return nil, fmt.Errorf("getting the previous version: %w", err)
	}
}

func ComposeView(db *db.Database, root db.ObjectiveId, vid db.VersionId) (any, error) {

}
