package sessions

import (
	"context"
	"fmt"
	"logbook/models/columns"
)

var ErrNoAuthentication = fmt.Errorf("no authentication")

func (a *Sessions) UserForSessionToken(ctx context.Context, token columns.SessionToken) (columns.UserId, error) {
	sid, err := a.q.SelectSessionByToken(ctx, token)
	if err != nil {
		return columns.ZeroUserId, fmt.Errorf("%w: selecting session id for session token from database: %w", ErrNoAuthentication, err)
	}
	return sid.Uid, nil
}
