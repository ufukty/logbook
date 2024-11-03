package endpoints

import (
	"logbook/cmd/objectives/database"
	"net/http"
)

type GetPlacementArrayRequest struct {
}

type GetPlacementArrayResponse struct {
	List []database.Objective
}

func (ep *Public) GetPlacementArray(w http.ResponseWriter, r *http.Request) {
	panic("to implement")
}
