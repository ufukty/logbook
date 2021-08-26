package database

import (
	"context"
	"errors"
)

func createDocument(document Document) (Document, bool, error) {
	query := `
		INSERT INTO "DOCUMENT" (
			"display_name",
			"user_id"
		)
		VALUES (
			$1, $2
		) 
		RETURNING
			"document_id", 
			"created_at"`
	err := pool.QueryRow(
		context.Background(),
		query,
		document.DisplayName,
		document.UserId,
	).Scan(
		&document.DocumentId,
		&document.CreatedAt,
	)
	if err != nil {
		return document, false, err
	}
	return document, true, nil
}

func getDocumentByDocumentId(documentId string) (Document, bool, error) {
	document := Document{DocumentId: documentId}
	query := `
		SELECT 
			"display_name",
			"user_id",
			"created_at"
		FROM
			"DOCUMENT"
		WHERE
			"document_id"=$1`
	err := pool.QueryRow(
		context.Background(),
		query,
		documentId,
	).Scan(
		&document.DisplayName,
		&document.UserId,
		&document.CreatedAt,
	)
	if err != nil {
		return document, false, err
	}
	return document, true, errors.New("")
}
