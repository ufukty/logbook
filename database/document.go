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

func (d *Document) Get() []error {
	query := `
		SELECT 
			"created_at",
			COALESCE("active_task", '00000000-0000-0000-0000-000000000000')
		FROM
			"DOCUMENT"
		WHERE
			"user_id"=$1
			AND "document_id"=$2`
	err := pool.QueryRow(
		context.Background(),
		query,
		d.UserId,
		d.DocumentId,
	).Scan(
		&d.CreatedAt,
		&d.ActiveTask,
	)
	if err != nil {
		return []error{err, ErrGetDocumentByDocumentId}
	}
	return nil
}

func (d *Document) GetHierarchicalPlacement() ([]string, []error) {
	taskIds := []string{}
	query := `
		SELECT * 
		FROM hierarchical_placement(
			v_user_id => $1,
			v_document_id => $2'
		)`
	rows, err := pool.Query(context.Background(), query, rd.userId, rd.documentId)
	if err != nil {
		return nil, []error{err, ErrGetHierarchicalViewItemsQuery}
	}
	for rows.Next() {
		var taskId string
		err := rows.Scan(&taskId)
		if err != nil {
			return nil, []error{err, ErrGetHierarchicalViewItemsScan}
		}
		taskIds = append(taskIds, taskId)
	}
	return taskIds, nil
}

func (d *Document) GetChronologicalPlacement() ([]string, []error) {
	taskIds := []string{}

	// The rows skipped by an OFFSET clause still have to be
	// computed inside the server; therefore a large OFFSET
	// might be inefficient.

	// So, hardcoded limit value is used for now. Later,
	// implement a mutable caching mechanism that doesn't
	// re-computes whole task ordering after each modification.
	// And instead, computes and overwrites changed range.
	query := `
		SELECT "task_id" 
		FROM "TASK"
		WHERE "user_id" = $1
			AND "document_id" = $2
			AND "archived" = FALSE
		ORDER BY "created_at" ASC
		LIMIT 100000
		`
	rows, err := pool.Query(context.Background(), query, rd.userId, rd.documentId)
	if err != nil {
		return nil, []error{err, ErrGetChronologicalViewItemsQuery}
	}
	for rows.Next() {
		var taskId string
		rows.Scan(&taskId)
		if err != nil {
			return nil, []error{err, ErrGetChronologicalViewItemsScan}
		}
		taskIds = append(taskIds, taskId)
	}
	return taskIds, nil
}
