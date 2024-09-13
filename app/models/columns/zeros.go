package columns

import "github.com/jackc/pgx/v5/pgtype"

const (
	ZeroAreaId          = AreaId("00000000-0000-0000-0000-000000000000")
	ZeroBookmarkId      = BookmarkId("00000000-0000-0000-0000-000000000000")
	ZeroCollaborationId = CollaborationId("00000000-0000-0000-0000-000000000000")
	ZeroGroupId         = GroupId("00000000-0000-0000-0000-000000000000")
	ZeroGroupInviteId   = GroupInviteId("00000000-0000-0000-0000-000000000000")
	ZeroLinkId          = LinkId("00000000-0000-0000-0000-000000000000")
	ZeroObjectId        = ObjectiveId("00000000-0000-0000-0000-000000000000")
	ZeroOperationId     = OperationId("00000000-0000-0000-0000-000000000000")
	ZeroPropertyId      = PropertiesId("00000000-0000-0000-0000-000000000000")
	ZeroUserId          = UserId("00000000-0000-0000-0000-000000000000")
	ZeroVersionId       = VersionId("00000000-0000-0000-0000-000000000000")
)

var (
	ZeroTimestamp = pgtype.Timestamp{}
)
