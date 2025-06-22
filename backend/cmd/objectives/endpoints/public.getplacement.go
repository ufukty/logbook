package endpoints

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/objectives/models/scalars"
	"logbook/cmd/sessions/endpoints"
	"logbook/internal/cookies"
	"logbook/internal/web/serialize"
	"logbook/models"
	"logbook/models/owners"
	"net/http"
)

type GetPlacementRequest struct {
	Root   models.Ovid             `route:"root"`
	Start  scalars.PlacementStart  `route:"start"`
	Length scalars.PlacementLength `route:"length"`
}

type GetPlacementResponse struct {
	Items []owners.DocumentItem `json:"items"`
}

// GET
func (p *Public) GetPlacement(w http.ResponseWriter, r *http.Request) {
	st, err := cookies.GetSessionToken(r)
	if err != nil {
		p.l.Println(fmt.Errorf("checking session token"))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	wi, err := p.sessions.WhoIs(&endpoints.WhoIsRequest{SessionToken: st})
	if err != nil {
		p.l.Println(fmt.Errorf("sessions.WhoIs: %w", err))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	bq := &GetPlacementRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			p.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
		return
	}

	dis, err := p.a.ViewBuilder(r.Context(), app.ViewBuilderParams{
		Viewer: wi.Uid,
		Root:   bq.Root,
		Start:  int(bq.Start),
		Length: int(bq.Length),
	})
	if err != nil {
		p.l.Println(fmt.Errorf("a.ViewBuilder: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bs := GetPlacementResponse{
		Items: dis,
	}
	if err := bs.Write(w); err != nil {
		p.l.Println(fmt.Errorf("writing json response: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
