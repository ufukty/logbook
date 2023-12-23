package database

import (
	"context"
)

func CreateDocument() (Document, []error) {
	document := Document{}
	query := `
		INSERT INTO "DOCUMENT" 
		DEFAULT VALUES
		RETURNING
			"document_id", 
			"created_at",
			COALESCE("active_task", '00000000-0000-0000-0000-000000000000')`
	err := pool.QueryRow(
		context.Background(),
		query,
	).Scan(
		&document.DocumentId,
		&document.CreatedAt,
		&document.ActiveTask,
	)
	if err != nil {
		return document, []error{err, ErrCreateDocument}
	}
	return document, nil
}

func GetDocumentByDocumentId(documentId string) (Document, []error) {
	document := Document{DocumentId: documentId}
	query := `
		SELECT 
			"created_at",
			COALESCE("active_task", '00000000-0000-0000-0000-000000000000')
		FROM
			"DOCUMENT"
		WHERE
			"document_id"=$1`
	err := pool.QueryRow(
		context.Background(),
		query,
		documentId,
	).Scan(
		&document.CreatedAt,
		&document.ActiveTask,
	)
	if err != nil {
		return document, []error{err, ErrGetDocumentByDocumentId}
	}
	return document, nil
}

func GetDocumentOverviewWithDocumentId(documentId string) ([]Task, []error) {
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
	rows, err := pool.Query(context.Background(), query, documentId)
	if err != nil {
		return nil, []error{err, ErrGetDocumentOverviewWithDocumentIdQuery}
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
			return nil, []error{err, ErrGetDocumentOverviewWithDocumentIdScan}
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetChronologicalViewItems(documentId string, limit int, offset int) ([]Task, []error) {
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
	rows, err := pool.Query(context.Background(), query, documentId, limit, offset)
	if err != nil {
		return nil, []error{err, ErrGetChronologicalViewItemsQuery}
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
			return nil, []error{err, ErrGetChronologicalViewItemsScan}
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
