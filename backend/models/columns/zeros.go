package columns

import "github.com/jackc/pgx/v5/pgtype"

const zeroUUID = "00000000-0000-0000-0000-000000000000"

const (
	ZeroAreaId            ControlAreaId     = zeroUUID
	ZeroBookmarkId        BookmarkId        = zeroUUID
	ZeroCollaborationId   CollaborationId   = zeroUUID
	ZeroCollaboratorId    CollaboratorId    = zeroUUID
	ZeroGroupId           GroupId           = zeroUUID
	ZeroGroupInviteId     GroupInviteId     = zeroUUID
	ZeroGroupMembershipId GroupMembershipId = zeroUUID
	ZeroLinkId            LinkId            = zeroUUID
	ZeroObjectiveId       ObjectiveId       = zeroUUID
	ZeroOperationId       OperationId       = zeroUUID
	ZeroPropertyId        PropertiesId      = zeroUUID
	ZeroUserId            UserId            = zeroUUID
	ZeroVersionId         VersionId         = zeroUUID
)

var (
	ZeroTimestamp = pgtype.Timestamp{}
)
