package app

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/models"
)

func (a *App) ListChildren(ctx context.Context, parent models.Ovid) ([]models.Ovid, error) {
	subs, err := a.oneshot.SelectSubLinks(ctx, database.SelectSubLinksParams{
		SupOid: parent.Oid,
		SupVid: parent.Vid,
	})
	if err != nil {
		return nil, fmt.Errorf("SelectSubLinks: %w", err)
	}
	children := []models.Ovid{}
	for _, sub := range subs {
		children = append(children, models.Ovid{Oid: sub.SubOid, Vid: sub.SubVid})
	}
	return children, nil
}
