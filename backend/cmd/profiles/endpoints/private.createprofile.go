package endpoints

import (
	"fmt"
	"logbook/cmd/profiles/app"
	"logbook/internal/web/serialize"

	"logbook/models/columns"
	"net/http"
)

type CreateProfileRequest struct {
	Uid       columns.UserId    `json:"uid"`
	Firstname columns.HumanName `json:"firstname"`
	Lastname  columns.HumanName `json:"lastname"`
}

// TODO: Authorization
// POST
func (p *Private) CreateProfile(w http.ResponseWriter, r *http.Request) {

	bq := &CreateProfileRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("binding request: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			p.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
		return
	}

	err := p.a.CreateProfile(r.Context(), app.CreateProfileParams{
		Uid:       bq.Uid,
		Firstname: bq.Firstname,
		Lastname:  bq.Lastname,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("app: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
