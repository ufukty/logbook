package private

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"net/http"
)

type MembershipCheckRequest struct {
	// TODO:
}

func (bq MembershipCheckRequest) validate() error {
	return validate.RequestFields(bq) // TODO: customize?
}

type MembershipCheckResponse struct {
	// TODO:
}

func (e *Endpoints) MembershipCheck(w http.ResponseWriter, r *http.Request) {
	bq := &MembershipCheckRequest{}

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

	bs := MembershipCheckResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
