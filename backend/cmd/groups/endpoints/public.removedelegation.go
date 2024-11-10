package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type RemoveDelegationRequest struct {
	SessionToken requests.Cookie[columns.SessionToken] `cookie:"sesion_token"`
	Delid        columns.DelegationId                  `json:"delid"`
}

type RemoveDelegationResponse struct {
	// TODO:
}

// POST
func (e *Public) RemoveDelegation(w http.ResponseWriter, r *http.Request) {
	bq := &RemoveDelegationRequest{}

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

	panic("to implement") // TODO:

	bs := RemoveDelegationResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
