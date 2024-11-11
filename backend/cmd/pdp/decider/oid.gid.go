package decider

import (
	"logbook/models/columns"
	"logbook/models/incoming"
)

func (d *Decider) OidGid(oid columns.ObjectiveId, gid columns.GroupId, act incoming.PolicyAction) error {
	return ErrUnderAuthorized
}
