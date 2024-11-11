package app

import (
	"context"
	"fmt"
	"logbook/models/columns"
)

var ErrCircularMembership = fmt.Errorf("two group directly or indirectly joining to each other is not possible")

// TODO: detect circular dependency
// TODO: enforce per-group member limit, 50?
// TODO: enforce per-user membership limit, 50?

func (a *App) InviteUserTypeMember(ctx context.Context, gid columns.GroupId, uid columns.UserId) (columns.GroupInviteId, error)
