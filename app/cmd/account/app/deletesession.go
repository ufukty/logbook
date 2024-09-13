package app

import (
	"context"
	"fmt"
	"logbook/models/columns"
)

func (a *App) DeleteSession(ctx context.Context, sid columns.SessionId) error {
	err := a.oneshot.DeleteSessionBySid(ctx, sid)
	if err != nil {
		return fmt.Errorf("deleting session in database: %w", err)
	}
	return nil
}
