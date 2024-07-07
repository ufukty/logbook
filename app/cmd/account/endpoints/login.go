package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/account/app"
	"logbook/cmd/account/app/average"
	"logbook/cmd/account/database"
	"logbook/internal/web/reqs"
	"logbook/internal/web/validate"
	"net/http"
)

type CreateSessionRequest struct {
	Email    database.Email `json:"email"`
	Password string         `json:"password"`
}

func (bq CreateSessionRequest) validate() error {
	return validate.All(map[string]validate.Validator{
		"email": bq.Email,
	})
}

func (e Endpoints) Login(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[CreateSessionRequest](r)
	if err != nil {
		e.l.Println(fmt.Errorf("binding: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		e.l.Println(fmt.Errorf("validation: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	session, err := e.app.Login(r.Context(), app.CreateSessionParameters{
		Email:    string(bq.Email),
		Password: string(bq.Password),
	})
	if err != nil {
		log.Println(fmt.Errorf("app.CreateSession: %w", err))
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
