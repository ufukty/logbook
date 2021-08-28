package document

import (
	"encoding/json"
	"fmt"
	"log"
	"logbook/main/database"
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	requestParameters := map[string]string{
		"ip-address": (*r).RemoteAddr,
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
	document, err = database.CreateDocument(database.Document{})

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
			"^ IP Address              : %s\n"+
			"^ User Agent              : %s", ipAddress, userAgent)
	log.Println(internalSuccessMessage)
}
