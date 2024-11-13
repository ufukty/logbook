package models

import (
	"fmt"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"strings"
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

func (ovid *Ovid) FromRoute(s string) error {
	if len(s) != len(columns.ZeroObjectId)+len(columns.ZeroVersionId)+1 {
		return fmt.Errorf("invalid length")
	}
	us := strings.Split(s, ":")
	if len(us) != 2 {
		return fmt.Errorf("invalid number of fragments")
	}
	ovid.Oid = columns.ObjectiveId(us[0])
	ovid.Vid = columns.VersionId(us[0])
	return nil
}

func (ovid *Ovid) ToRoute() (string, error) {
	return fmt.Sprintf("%s:%s", ovid.Oid, ovid.Vid), nil
}
