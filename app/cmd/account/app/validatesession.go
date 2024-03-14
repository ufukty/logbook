package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/app/average"
	"logbook/cmd/account/database"
	"time"
)

var ErrExpiredSession = fmt.Errorf("session is expired")

	session, err := a.queries.SelectSessionByToken(ctx, token)
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
