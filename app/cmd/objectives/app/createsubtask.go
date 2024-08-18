package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
	"slices"
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
	activepath, err := a.listActivePathToRock(ctx, params.Parent)
	if err == ErrLeftBehind {
		return ErrLeftBehind
	} else if err != nil {
		return fmt.Errorf("listActivePathToRock: %w", err)
	}

	op, err := a.queries.InsertOperation(ctx, database.InsertOperationParams{
		Subjectoid: params.Parent.Oid,
		Subjectvid: params.Parent.Vid,
		Actor:      params.Creator,
		OpType:     database.OpTypeObjCreateSubtask,
		OpStatus:   database.OpStatusAccepted,
	})
	if err != nil {
		return fmt.Errorf("InsertOperation: %w", err)
	}

	_, err = a.queries.InsertOpObjCreateSubtask(ctx, database.InsertOpObjCreateSubtaskParams{
		Opid:    op.Opid,
		Content: params.Content,
	})
	if err != nil {
		return fmt.Errorf("InsertOpObjCreateSubtask: %w", err)
	}

	props, err := a.queries.InsertProperties(ctx, database.InsertPropertiesParams{
		Content:   params.Content,
		Completed: false,
		Creator:   params.Creator,
		Owner:     params.Creator,
	})
	if err != nil {
		return fmt.Errorf("InsertProperties row: %w", err)
	}

	bup, err := a.queries.InsertBottomUpProps(ctx, database.InsertBottomUpPropsParams{
		SubtreeSize:      0,
		SubtreeCompleted: 0,
	})
	if err != nil {
		return fmt.Errorf("InsertBottomUpProps: %w", err)
	}

	obj, err := a.queries.InsertNewObjective(ctx, database.InsertNewObjectiveParams{
		CreatedBy: op.Opid,
		Pid:       props.Pid,
		Bupid:     bup.Bupid,
	})
	if err != nil {
		return fmt.Errorf("InsertNewObjective: %w", err)
	}

	_, err = a.queries.InsertActiveVidForObjective(ctx, database.InsertActiveVidForObjectiveParams{
		Oid: obj.Oid,
		Vid: obj.Vid,
	})
	if err != nil {
		return fmt.Errorf("InsertActiveVidForObjective: %w", err)
	}

	_, err = a.bubblink(ctx, slices.Insert(activepath, 0, models.Ovid{obj.Oid, obj.Vid}), op, bubblinkDeltaValues{SubtreeSize: 1})
	if err != nil {
		return fmt.Errorf("bubblink: %w", err)
	}

	return nil
}
