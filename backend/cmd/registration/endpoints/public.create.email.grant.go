package endpoints

import (
	"fmt"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
)

type CreateEmailGrantRequest struct {
	AntiCsrfToken transports.AntiCsrfToken `json:"acsrft"`
	Email         columns.Email            `json:"email"`
}

type CreateEmailGrantResponse struct {
	EmailGrant transports.EmailGrant `json:"email-grant"`
}

func (p *Public) CreateEmailGrant(w http.ResponseWriter, r *http.Request) {
	bq := &CreateEmailGrantRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		p.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	g, err := p.a.GrantEmail(r.Context(), bq.Email)
	if err != nil {
		p.l.Println(fmt.Errorf("a.GrantEmail: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreateEmailGrantResponse{
		EmailGrant: g,
	}
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
