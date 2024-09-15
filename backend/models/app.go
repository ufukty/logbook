package models

import (
	"fmt"
	"logbook/internal/web/validate"
	"logbook/models/columns"
)

// ObjectiveVersionedId: use to describe specific version of an objective
type Ovid struct {
	Oid columns.ObjectiveId `json:"oid"`
	Vid columns.VersionId   `json:"vid"`
}

var ZeroOvid = Ovid{Oid: columns.ZeroObjectId, Vid: columns.ZeroVersionId}

func (ovid Ovid) String() string {
	return fmt.Sprintf("(Oid: %q, Vid: %q)", ovid.Oid, ovid.Vid)
}

func (ovid Ovid) Validate() error {
	return validate.All(map[string]validate.Validator{
		"oid": ovid.Oid,
		"vid": ovid.Vid,
	})
}
