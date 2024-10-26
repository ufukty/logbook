package public

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"net/http"
)

type RemoveDelegationRequest struct {
	// TODO:
}

func (bq RemoveDelegationRequest) validate() error {
	return validate.RequestFields(bq) // TODO: customize?
}

type RemoveDelegationResponse struct {
	// TODO:
}

func (e *Endpoints) RemoveDelegation(w http.ResponseWriter, r *http.Request) {
	bq := &RemoveDelegationRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		e.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
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
