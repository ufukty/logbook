package app

import (
	"context"
	"fmt"
	"logbook/models/columns"
)

type WhoIsParams struct {
	SessionToken columns.SessionToken
}

func (a *App) WhoIs(ctx context.Context, params WhoIsParams) (columns.UserId, error) {
	ss, err := a.oneshot.SelectSessionByToken(ctx, params.SessionToken)
	if err != nil {
		return columns.ZeroUserId, fmt.Errorf("oneshot.SelectSessionByToken: %w", err)
	}

	if ss.Uid == columns.ZeroUserId {
		return columns.ZeroUserId, fmt.Errorf("unexpected zero user id found in database")
	}

	return ss.Uid, nil
}
