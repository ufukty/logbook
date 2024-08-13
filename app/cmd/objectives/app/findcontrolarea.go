package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"

	"github.com/jackc/pgx/v5"
)

var ErrControlAreaNotFound = fmt.Errorf("control area not found")

func (a *App) findControlArea(ctx context.Context, subject models.Ovid) (*database.ControlArea, error) {
	ap, err := a.listActivePathToRock(ctx, subject)
	if err != nil {
		return nil, fmt.Errorf("listActivePathToRock: %w", err)
	}
	for _, ascendant := range ap {
		ca, err := a.queries.SelectControlAreaOnObjective(ctx, ascendant.Oid)
		if err == pgx.ErrNoRows {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("SelectControlAreaOnObjective: %w", err)
		} else {
			return &ca, nil
		}
	}
	return nil, ErrControlAreaNotFound
}
