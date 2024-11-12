package app

import (
	"context"
	"logbook/models/columns"
	"logbook/models/transports"
)

type RespondToGroupInviteParams struct {
	Actor      columns.UserId
	Behalf     columns.GroupId
	Ginvid     columns.GroupInviteId // FIXME: differentiate
	MemberType transports.MemberType
	Response   transports.InviteResponse
}

func (a *App) RespondToGroupInvite(ctx context.Context, params RespondToGroupInviteParams) error
