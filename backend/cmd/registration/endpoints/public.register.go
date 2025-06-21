package endpoints

import (
	"fmt"
	"logbook/cmd/registration/app"
	"logbook/internal/web/serialize"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
)

type CreateAccountRequest struct {
	AntiCsrfToken transports.AntiCsrfToken `json:"acsrft"`

	Firstname columns.HumanName        `json:"firstname"`
	Lastname  columns.HumanName        `json:"lastname"`
	Birthday  transports.HumanBirthday `json:"birthday"`
	Country   transports.Country       `json:"country"`

	EmailGrant    transports.EmailGrant    `json:"email-grant"`
	PhoneGrant    transports.PhoneGrant    `json:"phone-grant"`
	PasswordGrant transports.PasswordGrant `json:"password-grant"`
}

// POST
func (p *Public) CreateAccount(w http.ResponseWriter, r *http.Request) {
	bq := &CreateAccountRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("binding: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			p.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
		return
	}

	err := p.a.Register(r.Context(), app.RegisterRequest{
		AntiCsrfToken: bq.AntiCsrfToken,
		Firstname:     bq.Firstname,
		Lastname:      bq.Lastname,
		Birthday:      bq.Birthday,
		Country:       bq.Country,
		EmailGrant:    bq.EmailGrant,
		PhoneGrant:    bq.PhoneGrant,
		PasswordGrant: bq.PasswordGrant,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("app.Register: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
