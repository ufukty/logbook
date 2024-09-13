package app

import (
	"context"
	"fmt"
	"logbook/models/columns"
)

func (a *App) Logout(ctx context.Context, token columns.SessionToken) error {
	err := a.oneshot.DeleteSessionByToken(ctx, token)
	if err != nil {
		// if err, ok := err.(*pgconn.PgError); ok {
		// 	return fmt.Errorf("")
		// }
		return fmt.Errorf("sending session deletion to database for session token %q: %w", token, err)
	}
	return nil
}
