package endpoints

import (
	"fmt"
	"logbook/cmd/account/api/public/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/reception"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type CreateProfileRequest struct {
	SessionToken requests.Cookie[columns.SessionToken] `cookie:"session_token"`
	Uid          columns.UserId                        `json:"uid"`
	Firstname    columns.HumanName                     `json:"firstname"`
	Lastname     columns.HumanName                     `json:"lastname"`
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
func (e Endpoints) CreateProfile(id reception.RequestId, store *reception.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &CreateProfileRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, redact(err), http.StatusInternalServerError)
		return fmt.Errorf("binding request: %w", err)
	}

	if err := bq.Validate(); err != nil {
		http.Error(w, redact(err), http.StatusBadRequest)
		return fmt.Errorf("validating the request parameters: %w", err)
	}

	err := e.a.CreateProfile(r.Context(), app.CreateProfileParams{
		SessionToken: bq.SessionToken.Value,
		Uid:          bq.Uid,
		Firstname:    bq.Firstname,
		Lastname:     bq.Lastname,
	})
	if err != nil {
		http.Error(w, redact(err), http.StatusInternalServerError)
		return fmt.Errorf("app: %w", err)
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
