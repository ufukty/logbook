package endpoints

import (
	"fmt"
	"logbook/cmd/account/api/public/app"
	"logbook/internal/average"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/registration/receptionist/decls"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type CreateSessionRequest struct {
	Email    columns.Email `json:"email"`
	Password string        `json:"password"`
}

func (bq CreateSessionRequest) validate() error {
	return validate.All(map[string]validate.Validator{
		"email": bq.Email,
	})
}

func (e Endpoints) Login(id decls.RequestId, store *decls.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &CreateSessionRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, redact(err), http.StatusBadRequest)
		return fmt.Errorf("binding: %w", err)
	}

	if err := bq.validate(); err != nil {
		http.Error(w, redact(err), http.StatusBadRequest)
		return fmt.Errorf("validation: %w", err)
	}

	session, err := e.a.Login(r.Context(), app.CreateSessionParameters{
		Email:    string(bq.Email),
		Password: string(bq.Password),
	})
	if err != nil {
		http.Error(w, redact(err), http.StatusInternalServerError)
		return fmt.Errorf("Login: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    string(session.Token),
		Expires:  session.CreatedAt.Time.Add(average.Week),
		HttpOnly: true,
		Secure:   true,
	})
	w.WriteHeader(http.StatusOK)

	return nil
}
