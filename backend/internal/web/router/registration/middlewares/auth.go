package middlewares

import (
	"logbook/internal/web/router/registration/decls"
	"net/http"
)

type auth struct {
}

func NewAuth() *auth {
	return &auth{}
}

func (auth) Handle(id decls.RequestId, store *decls.Store, w http.ResponseWriter, r *http.Request) error {
	return nil
}
