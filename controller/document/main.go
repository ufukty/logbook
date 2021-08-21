package document

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Document / get")
}

func Post(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Document / post")
}

func Details(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Document / details: %s\n", mux.Vars(r)["document_id"])
}

func Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Document / delete")
}
