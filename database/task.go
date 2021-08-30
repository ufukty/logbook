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

func GetTaskByParentId(parentId string) ([]Task, error) {
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
			"parent_id"=$1`
	rows, err := pool.Query(
		context.Background(),
		query,
		parentId,
	)
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

func UpdateTaskItem(task Task) (Task, error) {
	query := `
		UPDATE 
			"TASK" 
		SET 
			"content"=$2,
			"created_at"=$3,
			"degree"=$4,
			"depth"=$5,
			"parent_id"=$6,
			"task_group_id"=$7,
			"task_status"=$8
		WHERE 
			"task_id"=$1
		RETURNING
			"content",
			"created_at",
			"degree",
			"depth",
			"parent_id",
			"task_group_id",
			"task_status"`
	err := pool.QueryRow(
		context.Background(),
		query,
		task.TaskId,
		task.Content,
		task.CreatedAt,
		task.Degree,
		task.Depth,
		task.ParentId,
		task.TaskGroupId,
		task.TaskStatus,
	).Scan(
		&task.Content,
		&task.CreatedAt,
		&task.Degree,
		&task.Depth,
		&task.ParentId,
		&task.TaskGroupId,
		&task.TaskStatus,
	)
	return task, err
}
