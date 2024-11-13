package app

import (
	"context"
	"fmt"
	"logbook/cmd/sessions/database"
	"logbook/models/columns"
)

func (a *App) ActiveLoginsForUser(ctx context.Context, uid columns.UserId) ([]database.Login, error) {
	logins, err := a.oneshot.SelectLoginsByUid(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("marking login information as deleted in database: %w", err)
	}
	return logins, nil
}

func (a *App) LatestLoginForEmail(ctx context.Context, email columns.Email) (database.Login, error) {
	login, err := a.oneshot.SelectLatestLoginByEmail(ctx, email)
	if err != nil {
		return database.Login{}, fmt.Errorf("marking login information as deleted in database: %w", err)
	}
	return login, nil
}
