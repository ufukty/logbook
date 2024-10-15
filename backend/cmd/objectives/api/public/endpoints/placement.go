package endpoints

import (
	"logbook/cmd/objectives/database"
	"logbook/internal/web/router/registration/receptionist/decls"
	"net/http"
)

type GetPlacementArrayRequest struct {
}

type GetPlacementArrayResponse struct {
	List []database.Objective
}

func (ep *Endpoints) GetPlacementArray(rid decls.RequestId, store *decls.Store, w http.ResponseWriter, r *http.Request) error {
	return nil
}
