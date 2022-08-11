package controllers

import (
	"net/http"

	. "logbook/main/database"
	parameters "logbook/main/parameters"
	responder "logbook/main/responder"
)

func TaskCreate(w http.ResponseWriter, r *http.Request) {
	params := parameters.TaskCreate{}
	if err := params.InputSanitizer(r); err != nil {
		responder.ErrorHandler(w, r, http.StatusBadRequest, "Check your parameters", err)
	}

	Db.First()

	// check auth

	// creation of task
	// database.CreateTask(database.Task{
	// 	RevisionId:            "00000000-0000-0000-0000-000000000000",
	// 	OriginalCreatorUserId: "00000000-0000-0000-0000-000000000000",
	// 	ResponsibleUserId:     "00000000-0000-0000-0000-000000000000",
	// 	Content:               "Lorem ipsum dolor sit amet",
	// 	Notes:                 "Consectetur adipiscing elit",
	// })

	// creation of ownership role in PERM
	// database.CreatePermission(database.TaskPermission{
	// 	UserId: "00000000-0000-0000-0000-000000000000",
	// 	Role: "Role.Ownership",
	// })

	// check existence of super task

	// create link in TASK_LINK table

	// check permissions between task and user

	// create NewOperation

	// trigger task-props calculation

	responder.SuccessHandler(w, r, params.Response)
}
