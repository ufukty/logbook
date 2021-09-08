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
			"active_task"`
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
			"active_task"
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
			"completed_at", 
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
