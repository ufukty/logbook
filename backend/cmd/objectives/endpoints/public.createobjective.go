package endpoints

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/sessions/endpoints"
	"logbook/internal/cookies"
	"logbook/internal/web/serialize"
	"logbook/models"
	"logbook/models/columns"
	"net/http"
)

type CreateObjectiveRequest struct {
	Parent  models.Ovid              `json:"parent"`
	Content columns.ObjectiveContent `json:"content"`
}

type CreateObjectiveResponse struct {
	Oid columns.ObjectiveId `json:"oid"`
}

// TODO: Check user input for script tags in order to prevent XSS attempts
// POST
func (p *Public) CreateObjective(w http.ResponseWriter, r *http.Request) {
	st, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	wi, err := p.sessions.WhoIs(&endpoints.WhoIsRequest{SessionToken: st})
	if err != nil {
		p.l.Println(fmt.Errorf("sessions.WhoIs: %w", err))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bq := &CreateObjectiveRequest{}

	if err := bq.Parse(r); err != nil {
		fmt.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			p.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
		return
	}

	obj, err := p.a.CreateSubtask(r.Context(), app.CreateSubtaskParams{
		Parent:  bq.Parent,
		Content: bq.Content,
		Creator: wi.Uid,
	})
	if err != nil {
		fmt.Println(fmt.Errorf("app.CreateSubtask: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreateObjectiveResponse{
		Oid: obj.Oid,
	}
	if err := bs.Write(w); err != nil {
		fmt.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
