package document

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

// type TaskGroup struct {
// 	// GroupId string
// 	// Tasks []Task
// }

// type Document struct {
// 	DocumentId string
// 	TaskGroups []TaskGroup
// }

// type Dashboard struct {
// 	UserId    string
// 	Documents []Document
// }

type Document struct {
	DocumentID  string    `json:"document_id"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}

type Dashboard struct {
	UserId    string     `json:"user_id"`
	Documents []Document `json:"documents"`
}

func List(w http.ResponseWriter, r *http.Request) {

	userId := "0842c266-af1b-41bc-b180-653ca42dff82"

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5432/testdatabase")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	dashboard := Dashboard{
		UserId: userId,
	}

	// Read documents from database for given user
	rows, err2 := conn.Query(
		context.Background(),
		"SELECT \"document_id\", \"display_name\", \"created_at\" FROM \"DOCUMENT\" WHERE \"user_id\"=$1",
		userId,
	)

	// Read documents to structs
	for rows.Next() {
		document := Document{}
		if err = rows.Scan(&document.DocumentID, &document.DisplayName, &document.CreatedAt); err != nil {
			log.Println("Error in Document/Get at rows.Scan() for user: ", userId)
		}
		dashboard.Documents = append(dashboard.Documents, document)
	}

	if err2 != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	json.NewEncoder(w).Encode(dashboard)
	log.Println("Request processed.")
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
