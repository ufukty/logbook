package decider

import (
	"logbook/models/columns"
	"logbook/models/transports"
)

func (d *Decider) OidUid(oid columns.ObjectiveId, uid columns.UserId, act transports.PolicyAction) error {
	return ErrUnderAuthorized
}
