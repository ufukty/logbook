package endpoints

import (
	"fmt"
	"logbook/cmd/account/api/public/app"
	"logbook/internal/average"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type LoginRequest struct {
	Email    columns.Email `json:"email"`
	Password string        `json:"password"`
}

func (bq LoginRequest) validate() error {
	return validate.All(map[string]validate.Validator{
		"email": bq.Email,
	})
}

func (e Endpoints) Login(w http.ResponseWriter, r *http.Request) {
	bq := &LoginRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		e.l.Println(fmt.Errorf("binding: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		e.l.Println(fmt.Errorf("validation: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	session, err := e.a.Login(r.Context(), app.CreateSessionParameters{
		Email:    string(bq.Email),
		Password: string(bq.Password),
	})
	if err != nil {
		e.l.Println(fmt.Errorf("Login: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    string(session.Token),
		Expires:  session.CreatedAt.Time.Add(average.Week),
		HttpOnly: true,
		Secure:   true,
	})
	w.WriteHeader(http.StatusOK)
}
