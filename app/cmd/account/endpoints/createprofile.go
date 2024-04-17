package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/account/app"
	"logbook/cmd/account/database"
	"logbook/internal/web/reqs"
	"logbook/internal/web/validate"
	"net/http"
)

type CreateProfileRequest struct {
	SessionToken database.SessionToken `cookie:"session_token"`
	Uid          database.UserId       `json:"uid"`
	Firstname    database.HumanName    `json:"firstname"`
	Lastname     database.HumanName    `json:"lastname"`
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
	bq, err := reqs.ParseRequest[CreateProfileRequest](r)
	if err != nil {
		log.Println(fmt.Errorf("binding request: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := bq.Validate(); err != nil {
		log.Println(fmt.Errorf("validating the request parameters: %w", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ep.app.CreateProfile(r.Context(), app.CreateProfileParams{
		Uid:       bq.Uid,
		Firstname: bq.Firstname,
		Lastname:  bq.Lastname,
	})
	if err != nil {
		log.Println(fmt.Errorf("app: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
