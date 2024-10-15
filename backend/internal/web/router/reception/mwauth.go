package reception

import (
	"net/http"
)

type auth struct {
}

func NewAuth() *auth {
	return &auth{}
}

func (auth) Handle(id RequestId, store *Store, w http.ResponseWriter, r *http.Request) error {
	return nil
}
