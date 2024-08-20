package endpoints

import (
	"logbook/cmd/objectives/queries"
	"net/http"
)

type GetPlacementArrayRequest struct {
}

type GetPlacementArrayResponse struct {
	List []queries.Objective
}

func (ep *Endpoints) GetPlacementArray(w http.ResponseWriter, r *http.Request) {

}
