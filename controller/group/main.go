package group

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Group / details: %s\n", mux.Vars(r)["document_id"])
}
