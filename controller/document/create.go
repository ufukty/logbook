package document

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
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
	r.ParseForm()

	// Get userId from authorization/session information
	userId := "0842c266-af1b-41bc-b180-653ca42dff82"

	documentName := r.Form.Get("display_name")
	ipAddress := (*r).RemoteAddr
	userAgent := (*r).Header.Get("User-Agent")

	// Get display_name from request body
	if documentName == "" {
		eventId := uuid.New().String()
		publicErrorMessage := fmt.Sprintf(
			"400 Bad request.\n"+
				"You can not have a document without specifing a name.\n"+
				"Event ID: %s", eventId)
		internalErrorMessage := fmt.Sprintf(
			"[ERROR] Document/Create\n"+
				"^ Error reason            : No name supplied for creating an empty document.\n"+
				"^ Event ID                : %s\n"+
				"^ User ID                 : %s\n"+
				"^ Requested Document Name : %s\n"+
				"^ IP Address              : %s\n"+
				"^ User Agent              : %s", eventId, userId, documentName, ipAddress, userAgent)
		http.Error(w, publicErrorMessage, http.StatusBadRequest)
		log.Println(internalErrorMessage)
		return
	}

	// create document table record
	// & get the id of document just created
	documentId, errorCreateDocument := createDocumentRecord(userId, documentName)
	if errorCreateDocument != nil {
		eventId := uuid.New().String()
		publicErrorMessage := fmt.Sprintf(
			"500 Internal Server Error.\n"+
				"Please check back soon.\n"+
				"Event ID: %s", eventId)
		internalErrorMessage := fmt.Sprintf(
			"[ERROR] Document/Create\n"+
				"^ Error reason            : createDocumentRecord() raised error when write into the database.\n"+
				"^ Event ID                : %s\n"+
				"^ User ID                 : %s\n"+
				"^ Requested Document Name : %s\n"+
				"^ Created Document ID     : %s\n"+
				"^ IP Address              : %s\n"+
				"^ User Agent              : %s\n"+
				"^ Error details           : %s", eventId, userId, documentName, documentId, ipAddress, userAgent, errorCreateDocument)
		http.Error(w, publicErrorMessage, http.StatusInternalServerError)
		log.Println(internalErrorMessage)
		return
	}

	// created necessary group records for that document
	err_createGroups := createGroupRecords(userId, documentId)
	if err_createGroups != nil {
		eventId := uuid.New().String()
		publicErrorMessage := fmt.Sprintf(
			"500 Internal Server Error.\n"+
				"Please check back soon.\n"+
				"Event ID: %s", eventId)
		internalErrorMessage := fmt.Sprintf(
			"[ERROR] Document/Create\n"+
				"^ Error reason            : createGroupRecords() raised error when write into the database.\n"+
				"^ Event ID                : %s\n"+
				"^ User ID                 : %s\n"+
				"^ Requested Document Name : %s\n"+
				"^ Created Document ID     : %s\n"+
				"^ IP Address              : %s\n"+
				"^ User Agent              : %s\n"+
				"^ Error details           : %s", eventId, userId, documentName, documentId, ipAddress, userAgent, err_createGroups)
		http.Error(w, publicErrorMessage, http.StatusInternalServerError)
		log.Println(internalErrorMessage)
		return
	}

	// Prepare return object
	document := Document{DocumentId: documentId}

	taskGroups, errorTaskGroups := getTaskGroups(documentId)
	if errorTaskGroups != nil {
		eventId := uuid.New().String()
		publicErrorMessage := fmt.Sprintf(
			"410 Document might be corrupted.\n"+
				"Event ID: %s", eventId)
		internalErrorMessage := fmt.Sprintf(
			"[ERROR] Document/Create\n"+
				"^ Error reason            : getTaskGroups() raised error when read the database.\n"+
				"^ Event ID                : %s\n"+
				"^ User ID                 : %s\n"+
				"^ Requested Document Name : %s\n"+
				"^ Created Document ID     : %s\n"+
				"^ IP Address              : %s\n"+
				"^ User Agent              : %s\n"+
				"^ Error details           : %s", eventId, userId, documentName, documentId, ipAddress, userAgent, errorTaskGroups)
		http.Error(w, publicErrorMessage, http.StatusGone)
		log.Println(internalErrorMessage)
		return
	}

	document.TaskGroups = taskGroups
	document.TotalTaskGroups = len(taskGroups)

	for _, taskGroup := range taskGroups {
		tasks, errorTasks := getTasks(taskGroup.GroupId)
		if errorTasks != nil {
			eventId := uuid.New().String()
			publicErrorMessage := fmt.Sprintf(
				"410 Document might be corrupted.\n"+
					"Event ID: %s", eventId)
			internalErrorMessage := fmt.Sprintf(
				"[ERROR] Document/Create\n"+
					"^ Error reason            : getTasks() raised error when read the database.\n"+
					"^ Event ID                : %s\n"+
					"^ User ID                 : %s\n"+
					"^ Requested Document Name : %s\n"+
					"^ Created Document ID     : %s\n"+
					"^ Task Group ID           : %s\n"+
					"^ Task Group Type         : %s\n"+
					"^ IP Address              : %s\n"+
					"^ User Agent              : %s\n"+
					"^ Error details           : %s", eventId, userId, documentName, documentId, taskGroup.GroupId, taskGroup.GroupType, ipAddress, userAgent, errorTasks)
			http.Error(w, publicErrorMessage, http.StatusGone)
			log.Println(internalErrorMessage)
			return
		}
		taskGroup.Tasks = tasks
		taskGroup.TotalTasks = len(tasks)
	}

	json.NewEncoder(w).Encode(document)
	internalSuccessMessage := fmt.Sprintf(
		"[OK] Document/Create\n"+
			"^ User ID                 : %s\n"+
			"^ Requested Document Name : %s\n"+
			"^ Created Document ID     : %s\n"+
			"^ IP Address              : %s\n"+
			"^ User Agent              : %s", userId, documentName, documentId, ipAddress, userAgent)
	log.Println(internalSuccessMessage)
}
