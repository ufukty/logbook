package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateSubtaskParams struct {
	Creator columns.UserId
	Parent  models.Ovid
	Content string
}

// TODO: check prileges on parent
// DONE: create operations
// TODO: trigger task-props calculation
// TODO: transaction-commit-rollback
// DONE: bubblink
// DONE: mark active version for promoted ascendants
func (a *App) CreateSubtask(ctx context.Context, params CreateSubtaskParams) error {
	activepath, err := a.ListActivePathToRock(ctx, params.Parent)
	if err == ErrLeftBehind {
		return ErrLeftBehind
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
		Content: pgtype.Text{String: params.Content, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("inserting subtask creation details: %w", err)
	}

	obj, err := a.queries.InsertNewObjective(ctx, database.InsertNewObjectiveParams{
		CreatedBy: op.Opid,
		Props:     nil,
	})
	if err != nil {
		return fmt.Errorf("inserting the objective: %w", err)
	}

	// TODO: trigger computing props (async?)

	err = a.bubblink(ctx, obj, op, activepath)
	if err != nil {
		return fmt.Errorf("promoting the version change to ascendants: %w", err)
	}

	return nil
}
