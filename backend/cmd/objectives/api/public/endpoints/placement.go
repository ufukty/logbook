package endpoints

import (
	"logbook/cmd/objectives/api/public/middlewares"
	"logbook/cmd/objectives/database"
	"logbook/internal/web/router/receptionist"
	"net/http"
)

type GetPlacementArrayRequest struct {
}

type GetPlacementArrayResponse struct {
	List []database.Objective
}

func (ep *Endpoints) GetPlacementArray(rid receptionist.RequestId, store *middlewares.Store, w http.ResponseWriter, r *http.Request) error {
	return nil
}
