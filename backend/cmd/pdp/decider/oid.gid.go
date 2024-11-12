package decider

import (
	"logbook/models/columns"
	"logbook/models/transports"
)

func (d *Decider) OidGid(oid columns.ObjectiveId, gid columns.GroupId, act transports.PolicyAction) error {
	return ErrUnderAuthorized
}
