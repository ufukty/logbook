package app

import (
	"context"
	"fmt"
	"logbook/cmd/groups/database"
	"logbook/models/columns"
)

type CreateGroupParams struct {
	Actor     columns.UserId
	GroupName columns.GroupName
}

func (a *App) CreateGroup(ctx context.Context, params CreateGroupParams) (columns.GroupId, error) {
	g, err := a.oneshot.InsertNewGroup(ctx, database.InsertNewGroupParams{
		Name:    string(params.GroupName),
		Creator: params.Actor,
	})
	if err != nil {
		return columns.ZeroGroupId, fmt.Errorf("InsertNewGroup: %w", err)
	}
	return g.Gid, nil
}
