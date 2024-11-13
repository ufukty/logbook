package endpoints

import (
	"fmt"

	"logbook/models/columns"
	"net/http"
)

type CreateUserResponse struct {
	Uid columns.UserId `json:"uid"`
}

// TODO: Authorization
// POST
func (p *Private) CreateUser(w http.ResponseWriter, r *http.Request) {
	uid, err := p.a.CreateUser(r.Context())
	if err != nil {
		p.l.Println(fmt.Errorf("app: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	bq := CreateUserResponse{
		Uid: uid,
	}

	err = bq.Write(w)
	if err != nil {
		p.l.Println(fmt.Errorf(": %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
