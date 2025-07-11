package endpoints

import (
	"fmt"
	"logbook/cmd/sessions/app"
	"logbook/internal/cookies"
	"logbook/internal/web/serialize"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
)

type LoginRequest struct {
	Email    columns.Email       `json:"email"`
	Password transports.Password `json:"password"`
}

// POST
func (p *Public) Login(w http.ResponseWriter, r *http.Request) {
	bq := &LoginRequest{}

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

	session, err := p.a.Login(r.Context(), app.CreateSessionParameters{
		Email:    bq.Email,
		Password: bq.Password,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("Login: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	cookies.SetSessionToken(w, session.Token, session.CreatedAt.Time)
	w.WriteHeader(http.StatusOK)
}
