package document

import (
	"context"
	"errors"
	"log"
	"net/http"
)

func createDocumentRecord(userId string, documentName string) (string, error) {
	query := `INSERT INTO "DOCUMENT"("user_id", "display_name") VALUES($1, $2) RETURNING "document_id"`

	createdDocumentsId := ""
	err := PGXPool.QueryRow(context.Background(), query, userId, documentName).Scan(&createdDocumentsId)

	return createdDocumentsId, err
}

func createGroupRecords(userId string, documentId string) error {
	query := `INSERT INTO "TASK_GROUP"("document_id", "group_type") VALUES($1, $2)`

	var err error
	for _, groupType := range []string{"active", "paused", "ready-to-start", "plan", "dropped"} {
		_, err = PGXPool.Query(context.Background(), query, documentId, groupType)
	}

	return err
}

func Create(w http.ResponseWriter, r *http.Request) {

	// Get userId from authorization/session information
	userId := "0842c266-af1b-41bc-b180-653ca42dff82"

	// Get display_name from request body
	r.ParseForm()
	documentName := r.Form.Get("display_name")
	if documentName == "" {
		log_error(w, "POST /document", "r.Form.Get('diplay_name')", userId, errors.New(""))
	}

	// create document table record
	// & get the id of document just created
	documentId, error_createDoc := createDocumentRecord(userId, documentName)
	if error_createDoc != nil {
		log_error(w, "POST /document", "createDocumentRecord", userId, error_createDoc)
		return
	}

	// created necessary group records for that document
	err_createGroups := createGroupRecords(userId, documentId)
	if err_createGroups != nil {
		log_error(w, "POST /document", "createDocumentRecord createGroups", userId, err_createGroups)
		return
	}

	// Details(w, r)

	log.Println("POST /document: New document created\n^ userId:\t", userId, "\n^ documentId:\t", documentId)
}
