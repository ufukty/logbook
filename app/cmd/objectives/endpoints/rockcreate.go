package endpoints

import (
	"fmt"
	"log"
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

	if err := requests.ParseRequest(w, r, bq); err != nil {
		log.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		log.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := e.app.RockCreate(r.Context(), bq.UserId)
	if err != nil {
		log.Println(fmt.Errorf("app.RockCreate: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
