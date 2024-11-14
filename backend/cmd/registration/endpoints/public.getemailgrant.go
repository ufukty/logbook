package endpoints

import (
	"fmt"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
)

type GetEmailGrantRequest struct {
	AntiCsrfToken transports.AntiCsrfToken `json:"acsrft"`
	Email         columns.Email            `json:"email"`
}

type GetEmailGrantResponse struct {
	EmailGrant transports.EmailGrant `json:"email-grant"`
}

func (p *Public) GetEmailGrant(w http.ResponseWriter, r *http.Request) {
	bq := &GetEmailGrantRequest{}

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

	g, err := p.a.GrantEmail(r.Context(), bq.Email)
	if err != nil {
		p.l.Println(fmt.Errorf("a.GrantEmail: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := GetEmailGrantResponse{
		EmailGrant: g,
	}
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
