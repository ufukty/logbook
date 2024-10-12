package endpoints

import (
	"logbook/cmd/objectives/database"
	"logbook/internal/web/router/pipelines"
	"logbook/internal/web/router/pipelines/middlewares"
	"net/http"
)

type GetPlacementArrayRequest struct {
}

type GetPlacementArrayResponse struct {
	List []database.Objective
}

func (ep *Endpoints) GetPlacementArray(rid pipelines.RequestId, store *middlewares.Store, w http.ResponseWriter, r *http.Request) error {
	return nil
}
