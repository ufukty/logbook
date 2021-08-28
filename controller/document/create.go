package document

import (
	"encoding/json"
	"fmt"
	"log"
	"logbook/main/database"
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	ipAddress := (*r).RemoteAddr
	userAgent := (*r).Header.Get("User-Agent")

	var (
		document   database.Document
		task_group database.TaskGroup
		err        error
	)

	data, _ := json.Marshal(r.Header)
	fmt.Println(string(data))

	// create document table record
	document, err = database.CreateDocument(database.Document{})
	if err != nil {
		errorHandler(err, r, w)
	}

	for _, groupType := range []database.TaskStatus{
		database.Active, database.Archive, database.Drawer,
		database.Paused, database.ReadyToStart,
	} {
		task_group, err = database.CreateTaskGroup(
			database.TaskGroup{
				DocumentId:    document.DocumentId,
				TaskGroupType: database.TaskStatus(groupType),
			})
		if err != nil {
			errorHandler(err, r, w)
		}
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
