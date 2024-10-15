package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/registration/receptionist/decls"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type RockCreateRequest struct {
	UserId columns.UserId `json:"uid"`
}

func (bq RockCreateRequest) validate() error {
	return validate.RequestFields(bq)
}

func (e *Endpoints) RockCreate(rid decls.RequestId, store *decls.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &RockCreateRequest{}

	err := requests.ParseRequest(w, r, bq)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("parsing request: %w", err)
	}

	err = bq.validate()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return fmt.Errorf("validating request parameters: %w", err)
	}

	err = e.a.RockCreate(r.Context(), bq.UserId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("app.RockCreate: %w", err)
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
