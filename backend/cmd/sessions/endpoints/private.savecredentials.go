package endpoints

import (
	"fmt"
	"logbook/cmd/sessions/app"
	"logbook/internal/web/serialize"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
)

type SaveCredentialsRequest struct {
	Uid      columns.UserId      `json:"uid"`
	Email    columns.Email       `json:"email"`
	Password transports.Password `json:"password"`
}

func (p *Private) SaveCredentials(w http.ResponseWriter, r *http.Request) {
	bq := &SaveCredentialsRequest{}

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

	err := p.a.SaveCredentials(r.Context(), app.SaveCredentialsRequest{
		Uid:      bq.Uid,
		Email:    bq.Email,
		Password: bq.Password,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("p.a.SaveCredentials: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
