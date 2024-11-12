package app

import (
	"context"
	"logbook/models/columns"
	"logbook/models/transports"
)

type RespondToUserInviteParams struct {
	Actor      columns.UserId
	Ginvid     columns.GroupInviteId // FIXME: differentiate
	MemberType transports.MemberType
	Response   transports.InviteResponse
}

func (a *App) RespondToUserInvite(ctx context.Context, params RespondToUserInviteParams) error
