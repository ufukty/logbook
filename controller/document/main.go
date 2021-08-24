package document

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

var PGXPool *pgxpool.Pool

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

func log_error(w http.ResponseWriter, venue string, action string, userId string, err error) {
	eventId := uuid.New().String()
	log.Printf("%s: [ERROR] @ %s\n> [eventId]: %s\n> [userId]: %s\n> [details]: %s\n", venue, action, eventId, userId, err)
	http.Error(w, fmt.Sprintf("Internal server error. Please check back later. Event ID: %s", eventId), http.StatusInternalServerError)
}

func List(w http.ResponseWriter, r *http.Request) {

	// Get userId from authorization/session information
	userId := "0842c266-af1b-41bc-b180-653ca42dff82"

	// Read documents from database for given user
	rows, err_query := PGXPool.Query(
		context.Background(),
		"SELECT \"document_id\", \"display_name\", \"created_at\" FROM \"DOCUMENT\" WHERE \"user_id\"=$1",
		userId,
	)

	if err_query != nil {
		log_error(w, "/document/list", "conn.Query() 1st checkpoint", userId, err_query)
		return
	}

	// Read documents to structs
	dashboard := Dashboard{
		UserId: userId,
	}
	for rows.Next() {

		document := Document{}
		err_scan := rows.Scan(&document.DocumentID, &document.DisplayName, &document.CreatedAt)
		if err_scan != nil {
			log_error(w, "/document/list", "rows.Scan()", userId, err_scan)
			return
		}
		dashboard.Documents = append(dashboard.Documents, document)

	}

	if err_query != nil {
		log_error(w, "/document/list", "conn.Query() 2nd checkpoint", userId, err_query)
		return
	}

	json.NewEncoder(w).Encode(dashboard)
	log.Println("/document/list: Request proccessed for userId: ", userId)
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
