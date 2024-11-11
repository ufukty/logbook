package decider

import (
	"logbook/models/columns"
	"logbook/models/incoming"
)

func (d *Decider) OidUid(oid columns.ObjectiveId, uid columns.UserId, act incoming.PolicyAction) error {
	return ErrUnderAuthorized
}
