package models

import (
	"fmt"
	"logbook/models/columns"
	"strings"
)

// ObjectiveVersionedId: use to describe specific version of an objective
type Ovid struct {
	Oid columns.ObjectiveId `json:"oid"`
	Vid columns.VersionId   `json:"vid"`
}

var ZeroOvid = Ovid{Oid: columns.ZeroObjectiveId, Vid: columns.ZeroVersionId}

func (ovid Ovid) String() string {
	return fmt.Sprintf("(Oid: %q, Vid: %q)", ovid.Oid, ovid.Vid)
}

func (ovid Ovid) Validate() any {
	issues := map[string]any{}
	if issue := ovid.Oid.Validate(); issue != nil {
		issues["oid"] = issue
	}
	if issue := ovid.Vid.Validate(); issue != nil {
		issues["vid"] = issue
	}
	if len(issues) > 0 {
		return issues
	}
	return nil
}

func (ovid *Ovid) FromRoute(s string) error {
	if len(s) != len(columns.ZeroObjectiveId)+len(columns.ZeroVersionId)+1 {
		return fmt.Errorf("invalid length")
	}
	us := strings.Split(s, ":")
	if len(us) != 2 {
		return fmt.Errorf("invalid number of fragments")
	}
	ovid.Oid = columns.ObjectiveId(us[0])
	ovid.Vid = columns.VersionId(us[1])
	return nil
}

func (ovid *Ovid) ToRoute() (string, error) {
	return fmt.Sprintf("%s:%s", ovid.Oid, ovid.Vid), nil
}
