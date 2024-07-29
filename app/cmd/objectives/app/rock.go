package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models/columns"
)

// TODO: generate version number based on zero-vid
// TODO: insert objective using the version number
// TODO: insert bookmark using oid-vid
func (a *App) RockCreate(ctx context.Context, uid columns.UserId) error {
	v, err := a.queries.InsertVersion(ctx, columns.ZeroVersionId)
	if err != nil {
		return fmt.Errorf("queries.InsertVersion: %w", err)
	}

	o, err := a.queries.InsertObjective(ctx, database.InsertObjectiveParams{
		Vid:     v.Vid,
		Based:   columns.ZeroVersionId,
		Content: "",
		Creator: uid,
	})
	if err != nil {
		return fmt.Errorf("queries.InsertObjective: %w", err)
	}

	_, err = a.queries.InsertRock(ctx, database.InsertRockParams{
		Uid: uid,
		Oid: o.Oid,
		Vid: v.Vid,
	})
	if err != nil {
		return fmt.Errorf("queries.InsertRock: %w", err)
	}

	return nil
}
