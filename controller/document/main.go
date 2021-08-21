package document

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your-password"
	dbname   = "testdatabase"
)

func Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Document / get %s", password)

	db, err := sql.Open(
		"postgres", 
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// json.NewEncoder(w).Encode(linearized_tasks_)
	// log.Println("Request processed.")
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
