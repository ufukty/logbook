package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
	"time"
)

var week = time.Hour * 24 * 7

var ErrExpiredSession = fmt.Errorf("session is expired")

func (a *App) ValidateSession(ctx context.Context, sid database.SessionId) error {
	session, err := a.queries.SelectSession(ctx, sid)
	if err != nil {
		return fmt.Errorf("checking database: %w", err)
	}
	if session.Deleted {
		return ErrExpiredSession
	}
	if time.Now().Sub(session.CreatedAt.Time) > week {
		return ErrExpiredSession
	}
	return nil
}
