package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
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

func (e *Endpoints) RockCreate(w http.ResponseWriter, r *http.Request) {
	bq := &RockCreateRequest{}

	err := requests.ParseRequest(w, r, bq)
	if err != nil {
		e.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = bq.validate()
	if err != nil {
		e.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = e.a.RockCreate(r.Context(), bq.UserId)
	if err != nil {
		e.l.Println(fmt.Errorf("app.RockCreate: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
