package middlewares

import (
	"logbook/internal/web/router/reception"
	"net/http"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (Auth) Handle(id reception.RequestId, store *Store, w http.ResponseWriter, r *http.Request) error {
	return nil
}
