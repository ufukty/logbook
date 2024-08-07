package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
)

type CreateSubtaskParams struct {
	Creator columns.UserId
	Parent  models.Ovid
	Content string
}

// TODO: check prileges on parent
// DONE: create operations
// DONE: props
// TODO: transaction-commit-rollback
// DONE: bubblink
// DONE: mark active version for promoted ascendants
func (a *App) CreateSubtask(ctx context.Context, params CreateSubtaskParams) error {
	activepath, err := a.ListActivePathToRock(ctx, params.Parent)
	if err == ErrLeftBehind {
		return fmt.Errorf("checking active path: %w", ErrLeftBehind)
	} else if err != nil {
		return fmt.Errorf("checking the parent %s if it is inside active path: %w", params.Parent, err)
	}

	op, err := a.queries.InsertOperation(ctx, database.InsertOperationParams{
		Subjectoid: params.Parent.Oid,
		Subjectvid: params.Parent.Vid,
		Actor:      params.Creator,
		OpType:     database.OpTypeObjCreateSubtask,
		OpStatus:   database.OpStatusAccepted,
	})
	if err != nil {
		return fmt.Errorf("inserting the creation operation: %w", err)
	}

	_, err = a.queries.InsertOpObjCreateSubtask(ctx, database.InsertOpObjCreateSubtaskParams{
		Opid:    op.Opid,
		Content: params.Content,
	})
	if err != nil {
		return fmt.Errorf("inserting subtask creation details: %w", err)
	}

	props, err := a.queries.InsertProperties(ctx, database.InsertPropertiesParams{
		Content: params.Content,
		Creator: params.Creator,
	})
	if err != nil {
		return fmt.Errorf("inserting properties row: %w", err)
	}

	obj, err := a.queries.InsertNewObjective(ctx, database.InsertNewObjectiveParams{
		CreatedBy: op.Opid,
		Props:     props.Propid,
	})
	if err != nil {
		return fmt.Errorf("inserting the objective: %w", err)
	}

	// TODO: trigger computing props (async?)

	_, err = a.queries.InsertActiveVidForObjective(ctx, database.InsertActiveVidForObjectiveParams{
		Oid: obj.Oid,
		Vid: obj.Vid,
	})
	if err != nil {
		return fmt.Errorf("inserting active version: %w", err)
	}

	_, err = a.bubblink(ctx, append(activepath, models.Ovid{obj.Oid, obj.Vid}), op)
	if err != nil {
		return fmt.Errorf("promoting the version change to ascendants: %w", err)
	}

	return nil
}
