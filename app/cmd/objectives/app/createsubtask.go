package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/queries"
	"logbook/models"
	"logbook/models/columns"
	"slices"
)

type CreateSubtaskParams struct {
	Creator columns.UserId
	Parent  models.Ovid
	Content string
}

// TODO: check privileges on parent
// DONE: create operations
// DONE: props
// DONE: transaction-commit-rollback
// DONE: bubblink
// DONE: mark active version for promoted ascendants
func (a *App) CreateSubtask(ctx context.Context, params CreateSubtaskParams) (columns.ObjectiveId, error) {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := queries.New(tx)

	activepath, err := a.listActivePathToRock(ctx, q, params.Parent)
	if err == ErrLeftBehind {
		return columns.ZeroObjectId, ErrLeftBehind
	} else if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("listActivePathToRock: %w", err)
	}

	op, err := q.InsertOperation(ctx, queries.InsertOperationParams{
		Subjectoid: params.Parent.Oid,
		Subjectvid: params.Parent.Vid,
		Actor:      params.Creator,
		OpType:     queries.OpTypeObjCreateSubtask,
		OpStatus:   queries.OpStatusAccepted,
	})
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("InsertOperation: %w", err)
	}

	_, err = q.InsertOpObjCreateSubtask(ctx, queries.InsertOpObjCreateSubtaskParams{
		Opid:    op.Opid,
		Content: params.Content,
	})
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("InsertOpObjCreateSubtask: %w", err)
	}

	props, err := q.InsertProperties(ctx, queries.InsertPropertiesParams{
		Content:   params.Content,
		Completed: false,
		Creator:   params.Creator,
		Owner:     params.Creator,
	})
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("InsertProperties row: %w", err)
	}

	bup, err := q.InsertBottomUpProps(ctx, queries.InsertBottomUpPropsParams{
		Children:         0,
		SubtreeSize:      0,
		SubtreeCompleted: 0,
	})
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("InsertBottomUpProps: %w", err)
	}

	obj, err := q.InsertNewObjective(ctx, queries.InsertNewObjectiveParams{
		CreatedBy: op.Opid,
		Pid:       props.Pid,
		Bupid:     bup.Bupid,
	})
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("InsertNewObjective: %w", err)
	}

	_, err = q.InsertActiveVidForObjective(ctx, queries.InsertActiveVidForObjectiveParams{
		Oid: obj.Oid,
		Vid: obj.Vid,
	})
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("InsertActiveVidForObjective: %w", err)
	}

	_, err = a.bubblink(ctx, q, slices.Insert(activepath, 0, models.Ovid{obj.Oid, obj.Vid}), op, bubblinkDeltaValues{Children: 1, SubtreeSize: 1})
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("bubblink: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return columns.ZeroObjectId, fmt.Errorf("commit: %w", err)
	}

	return obj.Oid, nil
}
