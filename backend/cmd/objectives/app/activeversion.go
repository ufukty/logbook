package app

import (
	"context"
	"fmt"
	"logbook/models/columns"
)

func (a *App) GetActiveVersion(ctx context.Context, subject columns.ObjectiveId) (columns.VersionId, error) {
	act, err := a.oneshot.SelectActive(ctx, subject)
	if err != nil {
		return columns.ZeroVersionId, fmt.Errorf("SelectActive: %w", err)
	}
	return act.Vid, nil
}
