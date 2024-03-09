package app

import (
	"context"
	"fmt"
	"logbook/cmd/account/database"
)

type DeleteLoginParameters struct {
	Lid database.LoginId
}

func (a *App) DeleteLoginByLid(ctx context.Context, params DeleteLoginParameters) error {
	err := a.queries.DeleteLoginByLid(ctx, params.Lid)
	if err != nil {
		return fmt.Errorf("marking login information as deleted in database: %w", err)
	}
	return nil
}
