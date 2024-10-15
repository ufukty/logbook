package reception

import (
	"fmt"
	"net/http"
)

func pong(rid RequestId, store *Store, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "pong")
	return nil
}
