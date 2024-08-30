package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/queries"
	"logbook/models"
	"logbook/models/columns"
	"logbook/models/owners"
)

type GetObjectiveHistoryParams struct {
	Subject               models.Ovid
	IncludeAdministrative bool
}

func (a *App) GetObjectiveHistory(ctx context.Context, params GetObjectiveHistoryParams) ([]owners.OperationHistoryItem, error) {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := queries.New(tx)

	cursor := params.Subject.Vid
	stack := []owners.OperationHistoryItem{}
	for cursor != columns.ZeroVersionId {
		obj, err := q.SelectObjective(ctx, queries.SelectObjectiveParams{
			Oid: params.Subject.Oid,
			Vid: cursor,
		})
		if err != nil {
			return nil, fmt.Errorf("SelectObjective: %w", err)
		}
		op, err := q.SelectOperation(ctx, obj.CreatedBy)
		if err != nil {
			return nil, fmt.Errorf("SelectOperation: %w", err)
		}
		if params.IncludeAdministrative || op.Actor != columns.ZeroUserId { // bubblink
			stack = append(stack, owners.OperationHistoryItem{
				Version:   cursor,
				Type:      op.OpType,
				CreatedBy: op.Actor,
				CreatedAt: obj.CreatedAt.Time,
			})
		}
		cursor = obj.Based
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return stack, nil
}
