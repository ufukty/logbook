// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
	"logbook/models/columns"
)

type Tag struct {
	Tid       interface{}
	Vid       columns.VersionId
	Text      string
	Uid       columns.UserId
	Deleted   bool
	CreatedAt pgtype.Timestamp
}

type Tagging struct {
	Tid       columns.TagId
	Oid       columns.ObjectiveId
	Vid       columns.VersionId
	Deleted   bool
	CreatedAt pgtype.Timestamp
}
