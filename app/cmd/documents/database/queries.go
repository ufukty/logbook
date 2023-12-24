package database

import (
	"context"
)

func (db *Database) CreateDocument() (Document, []error) {
	document := Document{}
	query := `
		INSERT INTO "DOCUMENT"
		DEFAULT VALUES
		RETURNING
			"document_id",
			"created_at",
			COALESCE("active_task", '00000000-0000-0000-0000-000000000000')`
	err := db.pool.QueryRow(
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

func (db *Database) GetDocumentByDocumentId(documentId string) (Document, []error) {
	document := Document{DocumentId: documentId}
	query := `
		SELECT
			"created_at",
			COALESCE("active_task", '00000000-0000-0000-0000-000000000000')
		FROM
			"DOCUMENT"
		WHERE
			"document_id"=$1`
	err := db.pool.QueryRow(
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
