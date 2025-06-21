package endpoints

import (
	"fmt"
	"logbook/internal/web/serialize"
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

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			e.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
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
