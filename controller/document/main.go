package document

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"logbook/main/database"
)

func Init() {
	initErrors()
}

func List(w http.ResponseWriter, r *http.Request) {

	// Get userId from authorization/session information
	userId := "c66bc967-db0e-4911-a375-fc830db01a8b"

	var err error

	_, err = uuid.Parse(userId)
	if err != nil {
		http.Error(w, "Invalid user-id", http.StatusUnauthorized)
	}

	_, err = database.GetUserByUserId(userId)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			http.Error(w, "Invalid user-id", http.StatusUnauthorized)
		default:
			log.Println(err)
			http.Error(w, "Internal server error. Try again soon.", http.StatusInternalServerError)
		}
		return
	}

	documents, err := database.GetDocumentsByUserId(userId)
	if err != nil {
		log.Println("errorrrrr :", err)
		return
	}

	json.NewEncoder(w).Encode(documents)
	log.Println("GET /document/list: Request proccessed for userId: ", userId)
}

// func Update(w http.ResponseWriter, r *http.Request) {

// 	// Get userId from authorization/session information
// 	userId := "0842c266-af1b-41bc-b180-653ca42dff82"

// 	r.ParseForm()
// 	newName := r.Form.Get("name")
// 	if newName == "" {
// 		log_error(w, "PATCH /document/{document_id}", "r.Form.Get('name')", userId, errors.New("")) // TODO: NOT 502 INTERNAL SERVER ERROR
// 		return
// 	}

// 	documentId := mux.Vars(r)["document_id"]

// 	// Send SQL query
// 	query := "UPDATE \"DOCUMENT\" SET display_name=$1 WHERE \"document_id\"=$2 AND \"user_id\"=$3;"
// 	_, err_query := PGXPool.Query(context.Background(), query, newName, documentId, userId)

// 	if err_query != nil {
// 		log_error(w, "PATCH /document/{document_id}", "pgxpool.Pool.Query()", userId, err_query)
// 		return
// 	}

// 	log.Println("PATCH /document/{document_id}: Request proccessed for userId: ", userId)
// }

// // FIXME: ADD CSRF TOKEN
// func Delete(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	newName := r.Form.Get("name")
// 	if newName == "" {
// 		log_error(w, "PATCH /document/{document_id}", "r.Form.Get('name')", userId, errors.New("")) // TODO: NOT 502 INTERNAL SERVER ERROR
// 		return
// 	}
// 	fmt.Fprintln(w, "Document / delete")
// }
