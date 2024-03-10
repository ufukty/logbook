package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
)

func (a *App) DeleteSession(ctx context.Context, sid database.SessionId) error {
	err := a.queries.DeleteSessionBySid(ctx, sid)
	if err != nil {
		return fmt.Errorf("deleting session in database: %w", err)
	}
	return nil
}
