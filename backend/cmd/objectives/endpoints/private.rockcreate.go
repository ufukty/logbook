package endpoints

import (
	"fmt"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type RockCreateRequest struct {
	UserId columns.UserId `json:"uid"`
}

// POST
func (e *Private) RockCreate(w http.ResponseWriter, r *http.Request) {
	bq := &RockCreateRequest{}

	if err := bq.Parse(r); err != nil {
		e.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := validate.RequestFields(bq); err != nil {
		e.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := e.a.RockCreate(r.Context(), bq.UserId)
	if err != nil {
		e.l.Println(fmt.Errorf("app.RockCreate: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
