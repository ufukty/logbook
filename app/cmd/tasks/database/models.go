package database

import (
	"github.com/jackc/pgtype"
)

type VersioningConfig struct {
	Oid       ObjectiveId
	First     VersionId
	Effective VersionId
}

type Version struct {
	Vid   VersionId
	Based VersionId
}

type Action struct {
	Aid ActionId
	// PreviousOpId sql.NullString

	Vid   VersionId
	Based VersionId

	Oid ObjectiveId
	Uid string

	Summary ActionSummary
	Status  ActionStatus

	Creation   pgtype.Date
	ArchivedAt pgtype.Date
}

// objective or goal
type Objective struct {
	Oid      ObjectiveId
	Vid      VersionId
	Based    VersionId
	Type     ObjectiveType
	Content  string
	Creator  UserId
	Creation pgtype.Date
}

type Link struct {
	Lid LinkId

	SupOid ObjectiveId
	SupVid VersionId
	SubOid ObjectiveId
	SubVid VersionId

	Creation pgtype.Date
}

// Computed properties and user preferences per item per user
type ObjectiveView struct {
	Oid           ObjectiveId
	Vid           VersionId
	Uid           UserId
	Degree        NonNegativeNumber
	Depth         NonNegativeNumber
	ReadyToPickUp bool
	Completion    float64
	Fold          bool
}

type Bookmark struct {
	UserId       string
	Oid          ObjectiveId
	DisplayName  string
	RootBookmark string
	Creation     pgtype.Date
	DeletedAt    string
}

type ObjectivePermission struct {
}
