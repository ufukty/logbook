package database

import "context"

func createTaskGroup(taskGroup TaskGroup) (TaskGroup, bool, error) {
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
	if err != nil {
		return taskGroup, false, err
	}
	return taskGroup, true, nil
}

func getTaskGroupByGroupId(taskGroupId string) (TaskGroup, bool, error) {
	taskGroup := TaskGroup{}
	query := `
		SELECT 
			"task_group_id",
			"task_group_type",
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
		&taskGroup.DocumentId,
		&taskGroup.CreatedAt,
	)
	if err != nil {
		return taskGroup, false, err
	}
	return taskGroup, true, nil
}

func getTaskGroupsByDocumentId(documentId string) ([]TaskGroup, bool, error) {
	taskGroups := []TaskGroup{}
	query := `
		SELECT 
			"task_group_id",
			"task_group_type",
			"document_id",
			"created_at"
		FROM 
			"TASK_GROUP"
		WHERE
			"document_id"=$1`
	rows, err := pool.Query(context.Background(), query, documentId)
	if err != nil {
		return taskGroups, false, err
	}
	for rows.Next() {
		taskGroup := TaskGroup{}
		err = rows.Scan(
			&taskGroup.TaskGroupId,
			&taskGroup.TaskGroupType,
			&taskGroup.DocumentId,
			&taskGroup.CreatedAt,
		)
		if err != nil {
			continue
		}
		taskGroups = append(taskGroups, taskGroup)
	}
	return taskGroups, true, nil
}
