package app

import (
	"context"
	"errors"
	"fmt"
	"logbook/cmd/tasks/database"

	"github.com/jackc/pgx/v4"
)

// FIXME: Don't return error on pgx returns NoData but continoue to next iteration on loop
func (a *App) ListVersioningConfigForAncestry(ctx context.Context, ancestry []Ovid) ([]database.VersioningConfig, error) {
	vancestry := []database.VersioningConfig{}
	for _, c := range ancestry {
		vc, err := a.queries.SelectVersioningConfig(ctx, c.Oid)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			return vancestry, fmt.Errorf("SelectVersioningConfig(%q): %w", c.Oid, err)
		}
		vancestry = append(vancestry, vc)
	}
	return vancestry, nil
}
