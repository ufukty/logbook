package document

// func getTaskGroups(documentId string) ([]TaskGroup, error) {
// 	var taskGroups []TaskGroup

// 	query := `SELECT task_group_id, group_type FROM "TASK_GROUP" WHERE document_id=$1`

// 	rows, errQuery := PGXPool.Query(context.Background(), query, documentId)
// 	if errQuery != nil {
// 		return taskGroups, errQuery
// 	}

// 	for rows.Next() {
// 		taskGroup := TaskGroup{}
// 		errScan := rows.Scan(&taskGroup.GroupId, &taskGroup.GroupType)
// 		if errScan != nil {
// 			return taskGroups, errScan
// 		}
// 		taskGroups = append(taskGroups, taskGroup)
// 	}

// 	return taskGroups, errQuery
// }

// func getTasks(taskGroupId string) ([]Task, error) {
// 	var tasks []Task

// 	query := `
// 		SELECT
// 			task_id,
// 			parent_id,
// 			content,
// 			task_status,
// 			degree,
// 			depth,
// 			created_at
// 		FROM
// 			"TASK"
// 		WHERE
// 			task_group_id=$1`

// 	rows, errQuery := PGXPool.Query(context.Background(), query, taskGroupId)
// 	if errQuery != nil {
// 		return tasks, errQuery
// 	}

// 	for rows.Next() {
// 		task := Task{}
// 		errScan := rows.Scan(&task.TaskId, &task.ParentId, &task.Content, &task.TaskStatus, &task.Degree, &task.Depth, &task.CreatedAt)
// 		if errScan != nil {
// 			return tasks, errScan
// 		}
// 		tasks = append(tasks, task)
// 	}

// 	return tasks, errQuery
// }

// func checkDocumentExists() bool {
// 	return true
// }

// func Details(w http.ResponseWriter, r *http.Request) {

// 	// Get userId from authorization/session information
// 	userId := "0842c266-af1b-41bc-b180-653ca42dff82"

// 	ipAddress := (*r).RemoteAddr
// 	userAgent := (*r).Header.Get("User-Agent")
// 	documentId := mux.Vars(r)["document_id"]

// 	if !checkDocumentExists() {
// 		eventId := uuid.New().String()
// 		publicErrorMessage := fmt.Sprintf(
// 			"404 Document is not available.\n"+
// 				"It may not be exists or you might not have the appropriate rights to access.\n"+
// 				"Event ID: %s", eventId)
// 		internalErrorMessage := fmt.Sprintf(
// 			"[ERROR] Document/Details\n"+
// 				"^ Error reason          : Given document id is not valid, document does not exist.\n"+
// 				"^ Event ID              : %s\n"+
// 				"^ User ID               : %s\n"+
// 				"^ Requested Document ID : %s\n"+
// 				"^ IP Address            : %s\n"+
// 				"^ User Agent            : %s", eventId, userId, documentId, ipAddress, userAgent)
// 		http.Error(w, publicErrorMessage, http.StatusNotFound)
// 		log.Println(internalErrorMessage)
// 		return
// 	}

// 	document := Document{DocumentId: documentId}

// 	taskGroups, errorTaskGroups := getTaskGroups(documentId)
// 	if errorTaskGroups != nil {
// 		eventId := uuid.New().String()
// 		publicErrorMessage := fmt.Sprintf(
// 			"410 Document might be corrupted.\n"+
// 				"Event ID: %s", eventId)
// 		internalErrorMessage := fmt.Sprintf(
// 			"[ERROR] Document/Details\n"+
// 				"^ Error reason          : getTaskGroups() raised error when read the database.\n"+
// 				"^ Event ID              : %s\n"+
// 				"^ User ID               : %s\n"+
// 				"^ Requested Document ID : %s\n"+
// 				"^ IP Address            : %s\n"+
// 				"^ User Agent            : %s\n"+
// 				"^ Error details         : %s", eventId, userId, documentId, ipAddress, userAgent, errorTaskGroups)
// 		http.Error(w, publicErrorMessage, http.StatusGone)
// 		log.Println(internalErrorMessage)
// 		return
// 	}

// 	document.TaskGroups = taskGroups
// 	document.TotalTaskGroups = len(taskGroups)

// 	for _, taskGroup := range taskGroups {
// 		tasks, errorTasks := getTasks(taskGroup.GroupId)
// 		if errorTasks != nil {
// 			eventId := uuid.New().String()
// 			publicErrorMessage := fmt.Sprintf(
// 				"410 Document might be corrupted.\n"+
// 					"Event ID: %s", eventId)
// 			internalErrorMessage := fmt.Sprintf(
// 				"[ERROR] Document/Details\n"+
// 					"^ Error reason          : getTasks() raised error when read the database.\n"+
// 					"^ Event ID              : %s\n"+
// 					"^ User ID               : %s\n"+
// 					"^ Requested Document ID : %s\n"+
// 					"^ Task Group ID         : %s\n"+
// 					"^ Task Group Type       : %s\n"+
// 					"^ IP Address            : %s\n"+
// 					"^ User Agent            : %s\n"+
// 					"^ Error details         : %s", eventId, userId, documentId, taskGroup.GroupId, taskGroup.GroupType, ipAddress, userAgent, errorTasks)
// 			http.Error(w, publicErrorMessage, http.StatusGone)
// 			log.Println(internalErrorMessage)
// 			return
// 		}
// 		taskGroup.Tasks = tasks
// 		taskGroup.TotalTasks = len(tasks)
// 	}

// 	json.NewEncoder(w).Encode(document)
// 	internalSuccessMessage := fmt.Sprintf(
// 		"[OK] Document/Details\n"+
// 			"^ User ID               : %s\n"+
// 			"^ Requested Document ID : %s\n"+
// 			"^ IP Address            : %s\n"+
// 			"^ User Agent            : %s", userId, documentId, ipAddress, userAgent)
// 	log.Println(internalSuccessMessage)
// }
