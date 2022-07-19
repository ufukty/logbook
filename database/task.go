package database

import (
	"context"
)

// Returns a list of updated items in addition to
// the created task in first item.
func CreateTask(task Task) ([]Task, []error) {
	query := `
		SELECT 
			"task_id", 
			"document_id", 
			"parent_id", 
			"content", 
			"degree", 
			"depth", 
			"created_at", 
			"completed_at", 
			"ready_to_pick_up"
		FROM
			create_task($1, $2, $3)`
	rows, err := pool.Query(
		context.Background(),
		query,
		&task.DocumentId,
		&task.Content,
		&task.ParentId,
	)
	if err != nil {
		return nil, []error{err, ErrCreateTask}
	}
	tasks := []Task{}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(
			&task.TaskId,
			&task.DocumentId,
			&task.ParentId,
			&task.Content,
			&task.Degree,
			&task.Depth,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.ReadyToPickUp,
		)
		if err != nil {
			return nil, []error{err}
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTaskByTaskId(taskId string) (Task, []error) {
	task := Task{}
	query := `
		SELECT 
			"document_id", 
			"parent_id", 
			"content", 
			"degree", 
			"depth", 
			"created_at", 
			"completed_at", 
			"ready_to_pick_up"
		FROM 
			"TASK" 
		WHERE 
			"task_id"=$1`
	err := pool.QueryRow(context.Background(), query, taskId).Scan(
		&task.DocumentId,
		&task.ParentId,
		&task.Content,
		&task.Degree,
		&task.Depth,
		&task.CreatedAt,
		&task.CompletedAt,
		&task.ReadyToPickUp,
	)
	if err != nil {
		return task, []error{err, ErrGetTaskByTaskId}
	}
	return task, nil
}

func GetTasksByDocumentId(documentId string) ([]Task, []error) {
	tasks := []Task{}
	query := `
		SELECT 
			"document_id", 
			"parent_id", 
			"content", 
			"degree", 
			"depth", 
			"created_at", 
			"completed_at", 
			"ready_to_pick_up"
		FROM 
			"TASK" 
		WHERE 
			"document_id"=$1`
	rows, err := pool.Query(context.Background(), query, documentId)
	if err != nil {
		return tasks, []error{err, ErrGetTasksByDocumentIdQuery}
	}
	for rows.Next() {
		task := Task{}
		err = rows.Scan(
			&task.DocumentId,
			&task.ParentId,
			&task.Content,
			&task.Degree,
			&task.Depth,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.ReadyToPickUp,
		)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	if err != nil {
		return tasks, []error{err, ErrGetTasksByDocumentIdScan}
	}
	return tasks, nil
}

func GetTaskByParentId(parentId string) ([]Task, []error) {
	tasks := []Task{}
	query := `
	SELECT 
			"task_id", 
			"document_id", 
			"parent_id", 
			"content", 
			"degree", 
			"depth", 
			"created_at", 
			"completed_at", 
			"ready_to_pick_up"
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
		return tasks, []error{err, ErrGetTaskByParentIdQuery}
	}
	for rows.Next() {
		task := Task{}
		err = rows.Scan(
			&task.TaskId,
			&task.DocumentId,
			&task.ParentId,
			&task.Content,
			&task.Degree,
			&task.Depth,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.ReadyToPickUp,
		)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	if err != nil {
		return tasks, []error{err, ErrGetTaskByParentIdScan}
	}
	return tasks, nil
}

// func UpdateTaskItem(task Task) (Task, []error) {
// 	query := `
// 		UPDATE
// 			"TASK"
// 		SET
// 			"content"=$2,
// 			"created_at"=$3,
// 			"degree"=$4,
// 			"depth"=$5,
// 			"parent_id"=$6,
// 		WHERE
// 			"task_id"=$1
// 		RETURNING
// 			"content",
// 			"created_at",
// 			"degree",
// 			"depth",
// 			"parent_id",
// 			"task_group_id",
// 			"task_status"`
// 	err := pool.QueryRow(
// 		context.Background(),
// 		query,
// 		task.TaskId,
// 		task.Content,
// 		task.CreatedAt,
// 		task.Degree,
// 		task.Depth,
// 		task.ParentId,
// 	).Scan(
// 		&task.Content,
// 		&task.CreatedAt,
// 		&task.Degree,
// 		&task.Depth,
// 		&task.ParentId,
// 	)
// 	if err != nil {
// 		return task, []error{err, ErrUpdateTaskItem}
// 	}
// 	return task, nil
// }

// Returns a list of updated items in addition to
// the marked task in first item.
func (t *Task) MarkATaskDone() ([]Task, []error) {
	query := `SELECT mark_a_task_done($1)`
	rows, err := pool.Query(context.Background(), query, t.TaskId)
	if err != nil {
		return nil, []error{err, ErrMarkATaskDone}
	}
	tasks := []Task{}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(
			&task.TaskId,
			&task.DocumentId,
			&task.ParentId,
			&task.Content,
			&task.Degree,
			&task.Depth, 
			&task.CreatedAt,
			&task.CompletedAt,
			&task.ReadyToPickUp,
		)
		if err != nil {
			return nil, []error{err}
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Returns list of changed resources. Use it for cache invalidation.
func (t *Task) Reattach(newParentId string) ([]Task, []error) {
	query := `
		SELECT 
			"task_id"
		FROM reattach_task(
			v_task_id => $1, 
			v_new_parent_id => $2
		)`
	rows, err := pool.Query(context.Background(), query, t.TaskId, newParentId)
	if err != nil {
		return nil, []error{err, ErrMarkATaskDone}
	}
	tasks := []Task{}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.TaskId)
		if err != nil {
			return nil, []error{err}
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
