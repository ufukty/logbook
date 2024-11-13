package app

import (
	"context"
	"fmt"
	"logbook/cmd/sessions/database"
	"logbook/internal/average"
	"logbook/models/columns"
	"time"
)

var ErrExpiredSession = fmt.Errorf("session is expired")

func hasSessionExpired(session database.SessionStandard) bool {
	return time.Since(session.CreatedAt.Time) > average.Week
}

func (a *App) ValidateSession(ctx context.Context, token columns.SessionToken) error {
	session, err := a.oneshot.SelectSessionByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("checking database: %w", err)
	}
	if session.Deleted || hasSessionExpired(session) {
		return ErrExpiredSession
	}
	return nil
}
