package app

import (
	"fmt"
	"logbook/cmd/objectives/database"
	"logbook/internal/web/validate"
)

type CreateObjectiveAction struct {
	Parent  Ovid
	Content string
	Creator database.UserId
}

// ObjectiveVersionedId: use to describe specific version of an objective
type Ovid struct {
	Oid database.ObjectiveId
	Vid database.VersionId
}

func (ovid Ovid) String() string {
	return fmt.Sprintf("(Oid: %q, Vid: %q)", ovid.Oid, ovid.Vid)
}

func (ovid Ovid) Validate() error {
	return validate.All(map[string]validate.Validator{
		"oid": ovid.Oid,
		"vid": ovid.Vid,
	})
}
