package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/columns"
	"logbook/models/owners"
)

type ObjectiveHistoryParams struct {
	Subject models.Ovid
}

func (a *App) GetObjectiveHistory(ctx context.Context, params ObjectiveHistoryParams) ([]owners.OperationHistoryItem, error) {
	cursor := params.Subject.Vid
	stack := []owners.OperationHistoryItem{}
	for cursor != columns.ZeroVersionId {
		obj, err := a.queries.SelectObjective(ctx, database.SelectObjectiveParams{
			Oid: params.Subject.Oid,
			Vid: cursor,
		})
		if err != nil {
			return nil, fmt.Errorf("SelectObjective: %w", err)
		}
		op, err := a.queries.SelectOperation(ctx, obj.CreatedBy)
		if err != nil {
			return nil, fmt.Errorf("SelectOperation: %w", err)
		}
		if op.Actor == columns.ZeroUserId { // bubblink
			continue
		}
		stack = append(stack, owners.OperationHistoryItem{
			Version:   cursor,
			Type:      op.OpType,
			CreatedBy: op.Actor,
			CreatedAt: obj.CreatedAt.Time,
		})
		cursor = obj.Based
	}
	return stack, nil
}
