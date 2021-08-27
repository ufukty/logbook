package database

import (
	"context"
)

func CreateTask(task Task) (Task, error) {
	if err := checkTaskGroupId(task.TaskGroupId); err != nil {
		return task, err
	}
	query := `
		INSERT INTO "TASK" (
			"task_group_id",
			"parent_id",
			"content",
			"task_status",
			"degree",
			"depth"
		)
		VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING
			"task_id", "created_at"`
	err := pool.QueryRow(
		context.Background(),
		query,
		&task.TaskGroupId,
		&task.ParentId,
		&task.Content,
		&task.TaskStatus,
		&task.Degree,
		&task.Depth,
	).Scan(&task.TaskId, &task.CreatedAt)
	return task, exportError(err)
}

func GetTaskByTaskId(taskId string) (Task, error) {
	task := Task{}
	query := `
		SELECT 
			"content", 
			"created_at", 
			"degree", 
			"depth", 
			"parent_id", 
			"task_group_id",
			"task_id", 
			"task_status" 
		FROM 
			"TASK" 
		WHERE 
			"task_id"=$1`
	err := pool.QueryRow(context.Background(), query, taskId).Scan(
		&task.Content,
		&task.CreatedAt,
		&task.Degree,
		&task.Depth,
		&task.ParentId,
		&task.TaskGroupId,
		&task.TaskId,
		&task.TaskStatus,
	)
	return task, exportError(err)
}

func GetTasksByTaskGroupId(taskGroupId string) ([]Task, error) {
	tasks := []Task{}
	query := `
		SELECT 
			"content", 
			"created_at", 
			"degree", 
			"depth", 
			"parent_id", 
			"task_group_id",
			"task_id", 
			"task_status" 
		FROM 
			"TASK" 
		WHERE 
			"task_group_id"=$1`
	rows, err := pool.Query(context.Background(), query, taskGroupId)
	if err != nil {
		return tasks, exportError(err)
	}
	for rows.Next() {
		task := Task{}
		err = rows.Scan(
			&task.Content,
			&task.CreatedAt,
			&task.Degree,
			&task.Depth,
			&task.ParentId,
			&task.TaskGroupId,
			&task.TaskId,
			&task.TaskStatus,
		)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
