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
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)
	q := database.New(tx)

	ap, err := a.listActivePathToRock(ctx, q, subject)
	if err != nil {
		return nil, fmt.Errorf("listActivePathToRock: %w", err)
	}
	for _, ascendant := range ap {
		ca, err := q.SelectControlAreaOnObjective(ctx, ascendant.Oid)
		if err == pgx.ErrNoRows {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("SelectControlAreaOnObjective: %w", err)
		} else {
			return &ca, nil
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return nil, ErrControlAreaNotFound
}
