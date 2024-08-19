package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
	"logbook/models/owners"
)

func (a *App) GetMergedProps(ctx context.Context, subject models.Ovid) (owners.ObjectiveMergedProps, error) {
	obj, err := a.queries.SelectObjective(ctx, database.SelectObjectiveParams{
		Oid: subject.Oid,
		Vid: subject.Vid,
	})
	if err != nil {
		return owners.ObjectiveMergedProps{}, fmt.Errorf("SelectObjective: %w", err)
	}
	props, err := a.queries.SelectProperties(ctx, obj.Pid)
	if err != nil {
		return owners.ObjectiveMergedProps{}, fmt.Errorf("SelectProperties: %w", err)
	}
	bups, err := a.queries.SelectBottomUpProps(ctx, obj.Bupid)
	if err != nil {
		return owners.ObjectiveMergedProps{}, fmt.Errorf("SelectBottomUpProps: %w", err)
	}
	return owners.ObjectiveMergedProps{
		Content:          props.Content,
		Completed:        props.Completed,
		SubtreeSize:      bups.SubtreeSize,
		SubtreeCompleted: bups.SubtreeCompleted,
		Creator:          props.Creator,
		Owner:            props.Owner,
		CreatedAt:        props.CreatedAt.Time,
	}, nil
}
