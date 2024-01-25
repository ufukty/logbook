package database

import "github.com/jackc/pgtype"

type Action struct {
	Aid OperationId
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
