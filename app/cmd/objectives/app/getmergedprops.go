package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/queries"
	"logbook/models"
	"logbook/models/owners"
)

func (a *App) GetMergedProps(ctx context.Context, subject models.Ovid) (owners.ObjectiveMergedProps, error) {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return owners.ObjectiveMergedProps{}, fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := queries.New(tx)

	obj, err := q.SelectObjective(ctx, queries.SelectObjectiveParams{
		Oid: subject.Oid,
		Vid: subject.Vid,
	})
	if err != nil {
		return owners.ObjectiveMergedProps{}, fmt.Errorf("SelectObjective: %w", err)
	}
	props, err := q.SelectProperties(ctx, obj.Pid)
	if err != nil {
		return owners.ObjectiveMergedProps{}, fmt.Errorf("SelectProperties: %w", err)
	}
	bups, err := q.SelectBottomUpProps(ctx, obj.Bupid)
	if err != nil {
		return owners.ObjectiveMergedProps{}, fmt.Errorf("SelectBottomUpProps: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return owners.ObjectiveMergedProps{}, fmt.Errorf("commit: %w", err)
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
