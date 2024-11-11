package app

import (
	"context"
	"logbook/models/columns"
	"logbook/models/incoming"
)

type RespondToInviteParams struct {
	Invid      columns.GroupInviteId
	MemberType incoming.MemberType
	Response   incoming.InviteResponse
}

func (a *App) RespondToInvite(ctx context.Context, params RespondToInviteParams) (bool, error)
