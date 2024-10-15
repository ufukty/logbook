package middlewares

import (
	"fmt"
	"logbook/internal/web/router/registration/receptionist/decls"
	"net/http"
)

func Pong(rid decls.RequestId, store *decls.Store, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "pong")
	return nil
}
