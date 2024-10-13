package middlewares

import (
	"logbook/internal/web/router/receptionist"
	"net/http"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (Auth) Handle(id receptionist.RequestId, store *Store, w http.ResponseWriter, r *http.Request) error {
	return nil
}
