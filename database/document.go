package database

import (
	"context"
)

// Document object should have UserID. Rest will be overwritten.
func DocumentCreate(userId string) (*Document, []error) {
	doc := Document{}
	query := `
		INSERT INTO "DOCUMENT" ("user_id")
		VALUES ($1)
		RETURNING
			"document_id", 
			"user_id",
			"created_at"`
	err := pool.QueryRow(
		context.Background(),
		query,
		userId,
	).Scan(
		&doc.DocumentId,
		&doc.UserId,
		&doc.CreatedAt,
	)
	if err != nil {
		return &doc, []error{err, ErrCreateDocument}
	}
	return &doc, nil
}

func DocumentGet(userId string, documentId string) (*Document, []error) {
	doc := Document{}
	query := `
		SELECT 
			"document_id", 
			"user_id",
			"created_at"
		FROM "DOCUMENT"
		WHERE
			"user_id"=$1
			AND "document_id"=$2`
	err := pool.QueryRow(
		context.Background(),
		query,
		userId,
		documentId,
	).Scan(
		&doc.DocumentId,
		&doc.UserId,
		&doc.CreatedAt,
	)
	if err != nil {
		return &doc, []error{err, ErrGetDocumentByDocumentId}
	}
	return &doc, nil
}

// return value is intended to be used for caching, and only
// the client-requested portion will be sent to client.
func DocumentPlacementHierarchicalGet(userId string, documentId string) ([]string, []error) {
	taskIds := []string{}
	query := `
		SELECT * 
		FROM hierarchical_placement(
			v_user_id => $1,
			v_document_id => $2'
		)`
	rows, err := pool.Query(context.Background(), query, userId, documentId)
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

// return value is intended to be used for caching, and only
// the client-requested portion will be sent to client.
func DocumentPlacementChronologicalGet(userId string, documentId string) ([]string, []error) {
	taskIds := []string{}

	// TODO:
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
	rows, err := pool.Query(context.Background(), query, userId, documentId)
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
