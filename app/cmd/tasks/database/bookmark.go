package database

import (
	"github.com/jackc/pgtype"
)

type Bookmark struct {
	UserId       string
	Oid          ObjectiveId
	DisplayName  string
	RootBookmark string
	Creation     pgtype.Date
	DeletedAt    string
}
