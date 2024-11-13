package app

import (
	"context"
	"fmt"
	"logbook/models/columns"
)

var ErrNoRock = fmt.Errorf("rock not found")

func (a *App) RockGet(ctx context.Context, uid columns.UserId) (columns.ObjectiveId, error) {
	bs, err := a.oneshot.SelectTheRockForUser(ctx, uid)
	if err != nil {
		return columns.ZeroObjectiveId, fmt.Errorf("%w: SelectRockForUser: %w", ErrNoRock, err)
	}
	return bs.Oid, nil
}
