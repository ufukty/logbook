package endpoints

import (
	"fmt"
	"logbook/internal/web/validate"
	"logbook/models/transports"
	"net/http"
)

type CreatePasswordGrantRequest struct {
	AntiCsrfToken transports.AntiCsrfToken `json:"acsrft"`
	Password      transports.Password      `json:"password"`
}

type CreatePasswordGrantResponse struct {
	PasswordGrant transports.PasswordGrant `json:"password-grant"`
}

func (p *Public) CreatePasswordGrant(w http.ResponseWriter, r *http.Request) {
	bq := &CreatePasswordGrantRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := validate.RequestFields(bq); err != nil {
		p.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	g, err := p.a.GrantPassword(r.Context(), bq.Password)
	if err != nil {
		p.l.Println(fmt.Errorf("a.GrantPassword: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreatePasswordGrantResponse{
		PasswordGrant: g,
	}
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
