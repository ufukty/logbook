package database

import "context"

func CreateTaskGroup(taskGroup TaskGroup) (TaskGroup, error) {
	if err := checkDocumentId(taskGroup.DocumentId); err != nil {
		return taskGroup, err
	}
	query := `
		INSERT INTO "TASK_GROUP" (
			"document_id",
			"task_group_type"
		)
		VALUES (
			$1, $2
		)
		RETURNING
			"task_group_id",
			"created_at"
	`
	err := pool.QueryRow(
		context.Background(),
		query,
		taskGroup.DocumentId,
		taskGroup.TaskGroupType,
	).Scan(
		&taskGroup.TaskGroupId,
		&taskGroup.CreatedAt,
	)
	return taskGroup, exportError(err)
}

func GetTaskGroupByTaskGroupId(taskGroupId string) (TaskGroup, error) {
	taskGroup := TaskGroup{}
	query := `
		SELECT 
			"task_group_id",
			"task_group_type",
			"total_tasks",
			"document_id",
			"created_at"
		FROM 
			"TASK_GROUP"
		WHERE
			"task_group_id"=$1`
	err := pool.QueryRow(
		context.Background(),
		query,
		taskGroupId,
	).Scan(
		&taskGroup.TaskGroupId,
		&taskGroup.TaskGroupType,
		&taskGroup.TotalTasks,
		&taskGroup.DocumentId,
		&taskGroup.CreatedAt,
	)
	return taskGroup, exportError(err)
}

func GetTaskGroupsByDocumentId(documentId string) ([]TaskGroup, error) {
	taskGroups := []TaskGroup{}
	query := `
		SELECT 
			"task_group_id",
			"task_group_type",
			"total_tasks",
			"document_id",
			"created_at"
		FROM 
			"TASK_GROUP"
		WHERE
			"document_id"=$1`
	rows, err := pool.Query(context.Background(), query, documentId)
	if err != nil {
		return taskGroups, exportError(err)
	}
	for rows.Next() {
		taskGroup := TaskGroup{}
		err = rows.Scan(
			&taskGroup.TaskGroupId,
			&taskGroup.TaskGroupType,
			&taskGroup.TotalTasks,
			&taskGroup.DocumentId,
			&taskGroup.CreatedAt,
		)
		if err != nil {
			continue
		}
		taskGroups = append(taskGroups, taskGroup)
	}
	return taskGroups, nil
}
