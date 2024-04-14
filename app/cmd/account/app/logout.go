package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
)

func (a *App) Logout(ctx context.Context, token database.SessionToken) error {
	err := a.queries.DeleteSessionByToken(ctx, token)
	if err != nil {
		// if err, ok := err.(*pgconn.PgError); ok {
		// 	return fmt.Errorf("")
		// }
		return fmt.Errorf("sending session deletion to database for session token %q: %w", token, err)
	}
	return nil
}
