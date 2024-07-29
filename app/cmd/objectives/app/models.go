package app

import (
	"fmt"
	"logbook/internal/web/validate"
	"logbook/models/columns"
)

type CreateObjectiveAction struct {
	Parent  Ovid
	Content string
	Creator columns.UserId
}

// ObjectiveVersionedId: use to describe specific version of an objective
type Ovid struct {
	Oid columns.ObjectiveId `json:"oid"`
	Vid columns.VersionId   `json:"vid"`
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
