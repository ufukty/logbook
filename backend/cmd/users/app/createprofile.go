package app

import (
	"context"
	"fmt"
	"logbook/models/columns"
)

func (a *App) CreateUser(ctx context.Context) (columns.UserId, error) {
	u, err := a.oneshot.InsertUser(ctx)
	if err != nil {
		return columns.ZeroUserId, fmt.Errorf("insert into db: %w", err)
	}
	return u.Uid, nil
}
