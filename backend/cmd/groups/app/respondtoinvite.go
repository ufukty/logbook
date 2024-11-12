package app

import (
	"context"
	"logbook/models/columns"
	"logbook/models/transports"
)

type RespondToInviteParams struct {
	Invid      columns.GroupInviteId
	MemberType transports.MemberType
	Response   transports.InviteResponse
}

func (a *App) RespondToInvite(ctx context.Context, params RespondToInviteParams) (bool, error)
