package reception

import "net/http"

func NotFound(rid RequestId, store *Store, w http.ResponseWriter, r *http.Request) error {
	http.NotFound(w, r)
	return nil
}
