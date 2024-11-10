package endpoints

import (
	"fmt"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	"logbook/models"
	"logbook/models/columns"
	"net/http"
)

// TODO:
type TagAssignRequest struct {
	SessionToken requests.Cookie[columns.SessionToken] `cookie:"session_token"`
	Subject      models.Ovid                           `json:"subject"`
	Tid          columns.TagId                         `json:"tid"`
}

type TagAssignResponse struct {
	// TODO:
}

// POST
func (e *Public) TagAssign(w http.ResponseWriter, r *http.Request) {
	bq := &TagAssignRequest{}

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

	bs := TagAssignResponse{} // TODO:
	if err := requests.WriteJsonResponse(bs, w); err != nil {
		e.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
