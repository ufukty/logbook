package database

import (
	"context"
)

// Can throw those errors:
//    - ErrNoResult
func CreateDocument(document Document) (Document, error) {
	query := `
		INSERT INTO "DOCUMENT" 
		DEFAULT VALUES
		RETURNING
			"document_id", 
			"created_at",
			"total_task_groups"`
	err := pool.QueryRow(
		context.Background(),
		query,
	).Scan(
		&document.DocumentId,
		&document.CreatedAt,
		&document.TotalTaskGroups,
	)
	return document, exportError(err)
}

func CreateDocumentWithTaskGroups(document Document) (Document, error) {
	query := `
		SELECT 
			"document_id", 
			"created_at", 
			"total_task_groups"
		FROM create_document_with_task_groups();`
	err := pool.QueryRow(
		context.Background(),
		query,
	).Scan(
		&document.DocumentId,
		&document.CreatedAt,
		&document.TotalTaskGroups,
	)
	return document, exportError(err)
}

func GetDocumentByDocumentId(documentId string) (Document, error) {
	document := Document{DocumentId: documentId}
	query := `
		SELECT 
			"created_at",
			"total_task_groups"
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
		&document.TotalTaskGroups,
	)
	if err != nil {
		return document, err
	}
	return document, nil
}
