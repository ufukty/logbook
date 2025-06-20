package endpoints

import (
	"fmt"
	"logbook/internal/web/serialize"
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

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			p.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
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
