package task

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

var PGXPool *pgxpool.Pool

func Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Task / details: %s %s\n", mux.Vars(r)["document_id"], mux.Vars(r)["task_id"])
}

func Update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Task / details: %s %s\n", mux.Vars(r)["document_id"], mux.Vars(r)["task_id"])
}

func Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Task / details: %s %s\n", mux.Vars(r)["document_id"], mux.Vars(r)["task_id"])
}
