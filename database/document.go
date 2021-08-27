package database

import (
	"context"
)

func CreateDocument(document Document) (Document, error) {
	if err := checkUserId(document.UserId); err != nil {
		return document, err
	}
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
	return document, exportError(err)
}

func GetDocumentByDocumentId(documentId string) (Document, error) {
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
		return document, err
	}
	return document, nil
}

func GetDocumentsByUserId(userId string) ([]Document, error) {
	documents := []Document{}
	query := `
		SELECT 
			"document_id",
			"display_name",
			"user_id",
			"created_at"
		FROM 
			"DOCUMENT"
		WHERE
			"user_id"=$1`
	rows, err := pool.Query(
		context.Background(),
		query,
		userId,
	)
	if err != nil {
		return documents, err
	}
	for rows.Next() {
		document := Document{}
		err = rows.Scan(
			&document.DocumentId,
			&document.DisplayName,
			&document.UserId,
			&document.CreatedAt,
		)
	}
	return documents, err
}
