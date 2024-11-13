package app

import (
	"context"
	"errors"
	"fmt"
	"logbook/models/columns"

	"github.com/jackc/pgx/v5"
)

var (
	ErrSessionNotFound = fmt.Errorf("session not found")
	ErrUserNotFound    = fmt.Errorf("user not found")
	ErrProfileNotFound = fmt.Errorf("profile not found")
)

func (a App) WhoAmI(ctx context.Context, token columns.SessionToken) (columns.UserId, error) {
	session, err := a.oneshot.SelectSessionByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return columns.ZeroUserId, ErrSessionNotFound
		} else {
			return columns.ZeroUserId, fmt.Errorf("fetch session details from database: %w", err)
		}
	}

	if session.Deleted {
		return columns.ZeroUserId, ErrSessionNotFound
	}

	if hasSessionExpired(session) {
		return columns.ZeroUserId, ErrExpiredSession
	}

	// user, err := a.oneshot.SelectUserByUserId(ctx, session.Uid)
	// if err != nil {
	// 	return columns.ZeroUserId, fmt.Errorf("fetch user details from database: %w", err)
	// }

	// if user.Deleted {
	// 	return columns.ZeroUserId, ErrUserNotFound
	// }

	return session.Uid, nil
}
