package endpoints

import (
	"logbook/cmd/objectives/database"
	"logbook/internal/web/router/reception"
	"net/http"
)

type GetPlacementArrayRequest struct {
}

type GetPlacementArrayResponse struct {
	List []database.Objective
}

func (ep *Endpoints) GetPlacementArray(rid reception.RequestId, store *reception.Store, w http.ResponseWriter, r *http.Request) error {
	return nil
}
