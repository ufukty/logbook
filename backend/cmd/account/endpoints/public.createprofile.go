package endpoints

import (
	"fmt"
	"logbook/cmd/account/app"

	"logbook/internal/cookies"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type CreateProfileRequest struct {
	Uid       columns.UserId    `json:"uid"`
	Firstname columns.HumanName `json:"firstname"`
	Lastname  columns.HumanName `json:"lastname"`
}

func (params CreateProfileRequest) Validate() error {
	return validate.All(map[string]validate.Validator{
		"uid":       params.Uid,
		"firstname": params.Firstname,
		"lastname":  params.Lastname,
	})
}

// TODO: Authorization
// POST
func (p *Public) CreateProfile(w http.ResponseWriter, r *http.Request) {
	st, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bq := &CreateProfileRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("binding request: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	if err := bq.Validate(); err != nil {
		p.l.Println(fmt.Errorf("validating the request parameters: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	err = p.a.CreateProfile(r.Context(), app.CreateProfileParams{
		SessionToken: st,
		Uid:          bq.Uid,
		Firstname:    bq.Firstname,
		Lastname:     bq.Lastname,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("app: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
