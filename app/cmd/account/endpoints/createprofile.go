package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/account/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	database "logbook/models/columns"
	"net/http"
)

type CreateProfileRequest struct {
	SessionToken requests.Cookie[database.SessionToken] `cookie:"session_token"`
	Uid          database.UserId                        `json:"uid"`
	Firstname    database.HumanName                     `json:"firstname"`
	Lastname     database.HumanName                     `json:"lastname"`
}

func (params CreateProfileRequest) Validate() error {
	return validate.All(map[string]validate.Validator{
		"session_token": params.SessionToken,
		"uid":           params.Uid,
		"firstname":     params.Firstname,
		"lastname":      params.Lastname,
	})
}

// TODO: Authorization
func (ep Endpoints) CreateProfile(w http.ResponseWriter, r *http.Request) {
	bq := &CreateProfileRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		log.Println(fmt.Errorf("binding request: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	if err := bq.Validate(); err != nil {
		log.Println(fmt.Errorf("validating the request parameters: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	err := ep.app.CreateProfile(r.Context(), app.CreateProfileParams{
		SessionToken: bq.SessionToken.Value,
		Uid:          bq.Uid,
		Firstname:    bq.Firstname,
		Lastname:     bq.Lastname,
	})
	if err != nil {
		log.Println(fmt.Errorf("app: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}
}
