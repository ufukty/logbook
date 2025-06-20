package endpoints

import (
	"fmt"
	"logbook/internal/web/serialize"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
)

type CreatePhoneGrantRequest struct {
	AntiCsrfToken transports.AntiCsrfToken `json:"acsrft"`
	Phone         columns.Phone            `json:"phone"`
}

type CreatePhoneGrantResponse struct {
	PhoneGrant transports.PhoneGrant `json:"phone-grant"`
}

func (p *Public) CreatePhoneGrant(w http.ResponseWriter, r *http.Request) {
	bq := &CreatePhoneGrantRequest{}

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

	g, err := p.a.GrantPhone(r.Context(), bq.Phone)
	if err != nil {
		p.l.Println(fmt.Errorf("a.GrantEmail: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := CreatePhoneGrantResponse{
		PhoneGrant: g,
	}
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
