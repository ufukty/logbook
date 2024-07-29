package app

import (
	"context"
	"fmt"
	"logbook/models/columns"
)

type DeleteLoginParameters struct {
	Lid columns.LoginId
}

func (a *App) DeleteLoginByLid(ctx context.Context, params DeleteLoginParameters) error {
	err := a.queries.DeleteLoginByLid(ctx, params.Lid)
	if err != nil {
		return fmt.Errorf("marking login information as deleted in database: %w", err)
	}
	return nil
}
