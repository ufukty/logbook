package document

import (
	"encoding/json"
	"fmt"
	"log"
	"logbook/main/database"
	"net/http"

	"github.com/google/uuid"
)

func Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Get userId from authorization/session information
	userId := "d0d026b8-487d-45df-b4b6-08f54ab09615"

	requestParameters := map[string]string{
		"user-agent":          (*r).Header.Get("User-Agent"),
		"ip-address":          (*r).RemoteAddr,
		"input-document-name": r.Form.Get("display_name"),
	}

	desiredDocumentName := r.Form.Get("display_name")
	ipAddress := (*r).RemoteAddr
	userAgent := (*r).Header.Get("User-Agent")

	var (
		document   database.Document
		task_group database.TaskGroup
		err        error
	)
	// data, _ := json.Marshal(r.Body)
	// fmt.Println(string(data))

	// Get display_name from request body
	if desiredDocumentName == "" {
		errorHandler(ErrEmptyDocumentName, nil, requestParameters, w)
		return
	}

	// create document table record
	document, err = database.CreateDocument(
		database.Document{
			DisplayName: desiredDocumentName,
			UserId:      userId,
		})

	if err != nil {
		switch err { // FIXME:
		case database.ErrNoResult:
			log.Println("eagleeee", err)
		default:
			log.Println("agileeee", err)
		}
	}

	for _, groupType := range []database.TaskStatus{
		database.Active, database.Archive, database.Drawer,
		database.Paused, database.ReadyToStart,
	} {
		task_group, _ = database.CreateTaskGroup(
			database.TaskGroup{
				DocumentId:    document.DocumentId,
				TaskGroupType: database.TaskStatus(groupType),
			})
		document.TaskGroups = append(document.TaskGroups, task_group)
		document.TotalTaskGroups += 1
	}

	json.NewEncoder(w).Encode(document)
	internalSuccessMessage := fmt.Sprintf(
		"[OK] Document/Create\n"+
			"^ User ID                 : %s\n"+
			// "^ Requested Document Name : %s\n"+
			"^ Created Document ID     : %s\n"+
			"^ IP Address              : %s\n"+
			"^ User Agent              : %s", userId, desiredDocumentName, ipAddress, userAgent)
	log.Println(internalSuccessMessage)
}
