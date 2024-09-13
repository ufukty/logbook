package app

import (
	"context"
	"errors"
	"fmt"
	"logbook/cmd/account/database"
	"logbook/models/columns"

	"github.com/jackc/pgx/v5"
)

var (
	ErrSessionNotFound = fmt.Errorf("session not found")
	ErrUserNotFound    = fmt.Errorf("user not found")
	ErrProfileNotFound = fmt.Errorf("profile not found")
)

func (a App) WhoAmI(ctx context.Context, token columns.SessionToken) (*database.Profile, error) {
	session, err := a.oneshot.SelectSessionByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		} else {
			return nil, fmt.Errorf("fetch session details from database: %w", err)
		}
	}

	if session.Deleted {
		return nil, ErrSessionNotFound
	}

	if hasSessionExpired(session) {
		return nil, ErrExpiredSession
	}

	user, err := a.oneshot.SelectUserByUserId(ctx, session.Uid)
	if err != nil {
		return nil, fmt.Errorf("fetch user details from database: %w", err)
	}

	if user.Deleted {
		return nil, ErrUserNotFound
	}

	profile, err := a.oneshot.SelectLatestProfileByUid(ctx, session.Uid)
	if err != nil {
		return nil, ErrProfileNotFound
	}

	return &profile, nil
}
