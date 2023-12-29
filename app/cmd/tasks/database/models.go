package database

import (
	"github.com/jackc/pgtype"
)

type Action struct {
	Aid ActionId `json:"opid"`
	// PreviousOpId sql.NullString `json:"previous_opid"`

	Vid     VersionId `json:"vid"`
	PrevVid VersionId `json:"previous_vid"`

	Oid ObjectiveId `json:"oid"`

	UserId string `json:"uid"`

	Summary ActionSummary `json:"summary"`
	Status  ActionStatus  `json:"status"`

	CreatedAt  pgtype.Date `json:"created_at"`
	ArchivedAt pgtype.Date `json:"archived_at"`
}

// objective or goal
type Objective struct {
	Oid      ObjectiveId `json:"objective_id"`
	ParentId ObjectiveId `json:"parent_id"`
	Vid      VersionId   `json:"vid"`

	Creator UserId `json:"creator"`
	// ResponsibleUserId string      `json:"responsible_uid"`

	Text string `json:"text"`

	CreatedAt   pgtype.Date `json:"created_at"`
	CompletedAt pgtype.Date `json:"completed_at"` // nullable type
	ArchivedAt  pgtype.Date `json:"archived_at"`
}

type Link struct {
	Lid LinkId    `json:"lid"`
	Vid VersionId `json:"vid"`

	SupOid ObjectiveId `json:"sup_oid"` // immutable
	SupVid VersionId   `json:"sup_vid"` // immutable
	SubOid ObjectiveId `json:"sub_oid"` // immutable
	SubVid VersionId   `json:"sub_vid"` // immutable

	Index     int         `json:"index"`
	Type      LinkType    `json:"type"`
	CreatedAt pgtype.Date `json:"created_at"`
}

// Computed properties and user preferences per item per user
type ObjectiveView struct {
	Oid           ObjectiveId
	Vid           VersionId
	Uid           UserId
	Degree        NonNegativeNumber `json:"degree"`
	Depth         NonNegativeNumber `json:"depth"`
	ReadyToPickUp bool              `json:"ready_to_pick_up"`
	Completion    float64           `json:"completion"`
	Fold          bool              `json:"fold"`
}

type Bookmark struct {
	UserId       string      `json:"uid"`
	Oid          ObjectiveId `json:"objective_id"`
	DisplayName  string      `json:"display_name"`
	RootBookmark string      `json:"root_bookmark"`
	CreatedAt    pgtype.Date `json:"created_at"`
	DeletedAt    string      `json:"deleted_at"`
}

type ObjectivePermission struct {
}
