package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/app/average"
	"logbook/cmd/account/database"
	"time"
)

var ErrExpiredSession = fmt.Errorf("session is expired")

func (a *App) ValidateSession(ctx context.Context, sid database.SessionId) error {
	session, err := a.queries.SelectSessionBySid(ctx, sid)
	if err != nil {
		return fmt.Errorf("checking database: %w", err)
	}
	if session.Deleted {
		return ErrExpiredSession
	}
	if time.Now().Sub(session.CreatedAt.Time) > average.Week {
		return ErrExpiredSession
	}
	return nil
}
