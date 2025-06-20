package endpoints

import (
	"fmt"
	"logbook/cmd/sessions/endpoints"
	"logbook/internal/cookies"
	"logbook/models/columns"
	"net/http"
)

type TagCreationRequest struct {
	Title columns.TagTitle `json:"title"`
}

type TagCreationResponse struct {
	Tid columns.TagId `json:"tid"`
}

// POST
func (p *Public) TagCreation(w http.ResponseWriter, r *http.Request) {
	st, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	_, err = p.sessions.WhoIs(&endpoints.WhoIsRequest{
		SessionToken: st,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("sessions.WhoIs: %w", err))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bq := &TagCreationRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		p.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	panic("to implement") // TODO:

	bs := TagCreationResponse{} // TODO:
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
