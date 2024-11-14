package endpoints

import (
	"fmt"
	"logbook/internal/web/validate"
	"logbook/models/transports"
	"net/http"
)

type GetPasswordGrantRequest struct {
	AntiCsrfToken transports.AntiCsrfToken `json:"acsrft"`
	Password      transports.Password      `json:"password"`
}

type GetPasswordGrantResponse struct {
	PasswordGrant transports.PasswordGrant `json:"password-grant"`
}

func (p *Public) GetPasswordGrant(w http.ResponseWriter, r *http.Request) {
	bq := &GetPasswordGrantRequest{}

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

	bs := GetPasswordGrantResponse{
		PasswordGrant: g,
	}
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
