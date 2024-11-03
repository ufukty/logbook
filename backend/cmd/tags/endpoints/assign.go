package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"net/http"
)

type TagAssignRequest struct {
	// TODO:
}

func (bq TagAssignRequest) validate() error {
	panic("to implement") // TODO:
	return nil
}

type TagAssignResponse struct {
	// TODO:
}

func (e *Endpoints) TagAssign(w http.ResponseWriter, r *http.Request) {
	bq := &TagAssignRequest{}

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

	bs := TagAssignResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
