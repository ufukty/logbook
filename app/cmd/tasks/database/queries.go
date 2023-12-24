package database

import (
	"context"
	"fmt"
)

// Returns a list of updated items in addition to
// the created task in first item.
func (db *Database) CreateTask(task *Task) ([]Task, error) {
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
	rows, err := db.pool.Query(
		context.Background(),
		query,
		&task.DocumentId,
		&task.Content,
		&task.ParentId,
	)
	if err != nil {
		return nil, fmt.Errorf("running the query: %w", err)
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
			return nil, fmt.Errorf("scanning query results: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *Database) GetTaskByTaskId(taskId string) (Task, error) {
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
	err := db.pool.QueryRow(context.Background(), query, taskId).Scan(
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
		return task, fmt.Errorf("running the query: %w", err)
	}
	return task, nil
}

func (db *Database) ListTasksInDocument(documentId string) ([]Task, error) {
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
	rows, err := db.pool.Query(context.Background(), query, documentId)
	if err != nil {
		return tasks, fmt.Errorf("running the query: %w", err)
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
		return nil, fmt.Errorf("scanning query results: %w", err)
	}
	return tasks, nil
}

func (db *Database) GetSubitems(iid string) ([]Task, error) {
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
	rows, err := db.pool.Query(
		context.Background(),
		query,
		iid,
	)
	if err != nil {
		return tasks, fmt.Errorf("running the query: %w", err)
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
		return nil, fmt.Errorf("scanning query results: %w", err)
	}
	return tasks, nil
}

// func (db *Database)UpdateTaskItem(task Task) (Task, error) {
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
// 	err := db.pool.QueryRow(
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
//      return nil, fmt.Errorf("scanning query results: %w", err)
// 	}
// 	return task, nil
// }

// Returns a list of updated items in addition to
// the marked task in first item.
func (db *Database) MarkDone(taskId string) ([]Task, error) {
	query := `SELECT mark_a_task_done($1)`
	rows, err := db.pool.Query(context.Background(), query, taskId)
	if err != nil {
		return nil, fmt.Errorf("running the query: %w", err)
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
			return nil, fmt.Errorf("scanning query results: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Returns a list of updated items in addition to
// the reattached task in first item.
func (db *Database) ReattachTask(taskId string, newParentId string) ([]Task, error) {
	query := `SELECT reattach_task($1, $2)`
	rows, err := db.pool.Query(context.Background(), query, taskId, newParentId)
	if err != nil {
		return nil, fmt.Errorf("running the query: %w", err)
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
			return nil, fmt.Errorf("scanning query results: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *Database) GetDocumentOverviewWithDocumentId(documentId string) ([]Task, error) {
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
			COALESCE("completed_at", '0001-01-01'),
			"ready_to_pick_up"
		FROM document_overview($1)`
	rows, err := db.pool.Query(context.Background(), query, documentId)
	if err != nil {
		return nil, fmt.Errorf("running the query: %w", err)
	}
	for rows.Next() {
		task := Task{}
		rows.Scan(
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
			return nil, fmt.Errorf("scanning query results: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *Database) GetChronologicalViewItems(documentId string, limit int, offset int) ([]Task, error) {
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
		FROM "TASK"
		WHERE "document_id" = $1
		ORDER BY "created_at" ASC
		LIMIT $2 OFFSET $3
		`
	rows, err := db.pool.Query(context.Background(), query, documentId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("running the query: %w", err)
	}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(
			&(task.TaskId),
			&(task.DocumentId),
			&(task.ParentId),
			&(task.Content),
			&(task.Degree),
			&(task.Depth),
			&(task.CreatedAt),
			&(task.CompletedAt),
			&(task.ReadyToPickUp),
		)
		if err != nil {
			return nil, fmt.Errorf("scanning query results: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
