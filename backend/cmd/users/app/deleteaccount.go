package app

import (
	"context"
	"fmt"
	"logbook/cmd/users/database"
	"logbook/models/columns"
)

func (a *App) DeleteAccount(ctx context.Context, uid columns.UserId) error {
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := database.New(tx)

	// sessions, err := q.SelectActiveSessionsByUid(ctx, uid)
	// if err != nil {
	// 	return fmt.Errorf("selecting active sessions for user %q from database: %w", uid, err)
	// }

	// for _, session := range sessions {
	// 	err := q.DeleteSessionBySid(ctx, session.Sid)
	// 	if err != nil {
	// 		return fmt.Errorf("sending deletion request of session %q to database: %w", session.Sid, err)
	// 	}
	// }

	err = q.DeleteUserByUid(ctx, uid)
	if err != nil {
		return fmt.Errorf("marking user deleted on database: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}
