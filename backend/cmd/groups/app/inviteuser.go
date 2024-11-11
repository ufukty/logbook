package app

import (
	"context"
	"logbook/models/columns"
)

func (a *App) InviteGroupTypeMember(ctx context.Context, group, member columns.GroupId) (columns.GroupInviteId, error)
