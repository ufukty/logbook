package endpoints

import (
	"fmt"
	"logbook/cmd/pdp/decider"
	"logbook/internal/web/serialize"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
)

type OidUidRequest struct {
	Actor    columns.UserId          `route:"uid"`
	Resource columns.ObjectiveId     `route:"oid"`
	Action   transports.PolicyAction `route:"action"`
}

// GET
func (p *Private) OidUid(w http.ResponseWriter, r *http.Request) {
	bq := &OidUidRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			p.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
		return
	}

	err := p.d.OidUid(bq.Resource, bq.Actor, bq.Action)
	if err != nil && err != decider.ErrUnderAuthorized {
		p.l.Println(fmt.Errorf("p.d.OidUid: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err == decider.ErrUnderAuthorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
