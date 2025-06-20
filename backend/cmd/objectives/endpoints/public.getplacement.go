package endpoints

import (
	"fmt"
	"logbook/cmd/objectives/app"
	"logbook/cmd/sessions/endpoints"
	"logbook/internal/cookies"
	"logbook/models"
	"logbook/models/owners"
	"net/http"
	"strconv"
)

type PlacementStart int

func (p *PlacementStart) FromRoute(src string) error {
	a, err := strconv.Atoi(src)
	if err != nil {
		return fmt.Errorf("atoi: %w", err)
	}
	*p = PlacementStart(a)
	return nil
}

func (p PlacementStart) ToRoute() (string, error) {
	return strconv.Itoa(int(p)), nil
}

func (p PlacementStart) Validate() error {
	if 0 <= p && p < 10000 {
		return fmt.Errorf("out of range")
	}
	return nil
}

type PlacementLength int

func (p *PlacementLength) FromRoute(src string) error {
	a, err := strconv.Atoi(src)
	if err != nil {
		return fmt.Errorf("atoi: %w", err)
	}
	*p = PlacementLength(a)
	return nil
}

func (p PlacementLength) ToRoute() (string, error) {
	return strconv.Itoa(int(p)), nil
}

func (p PlacementLength) Validate() error {
	if 0 <= p && p < 10000 {
		return fmt.Errorf("out of range")
	}
	return nil
}

type GetPlacementRequest struct {
	Root   models.Ovid     `route:"root"`
	Start  PlacementStart  `route:"start"`
	Length PlacementLength `route:"length"`
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
		p.l.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
