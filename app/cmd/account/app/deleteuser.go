package app

import (
	"context"
	"fmt"
	database "logbook/models/columns"
)

func (a *App) DeleteUser(ctx context.Context, uid database.UserId) error {
	sessions, err := a.oneshot.SelectActiveSessionsByUid(ctx, uid)
	if err != nil {
		return fmt.Errorf("selecting active sessions for user %q from database: %w", uid, err)
	}

	for _, session := range sessions {
		err := a.oneshot.DeleteSessionBySid(ctx, session.Sid)
		if err != nil {
			return fmt.Errorf("sending deletion request of session %q to database: %w", session.Sid, err)
		}
	}

	err = a.oneshot.DeleteUserByUid(ctx, uid)
	if err != nil {
		return fmt.Errorf("marking user deleted on database: %w", err)
	}

	return nil
}
