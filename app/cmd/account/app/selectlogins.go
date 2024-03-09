package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
)

func (a *App) ActiveLoginsForUser(ctx context.Context, uid database.UserId) ([]database.Login, error) {
	logins, err := a.queries.SelectLoginsByUid(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("marking login information as deleted in database: %w", err)
	}
	return logins, nil
}

func (a *App) LatestLoginForEmail(ctx context.Context, email database.Email) (database.Login, error) {
	login, err := a.queries.SelectLatestLoginByEmail(ctx, string(email))
	if err != nil {
		return database.Login{}, fmt.Errorf("marking login information as deleted in database: %w", err)
	}
	return login, nil
}
