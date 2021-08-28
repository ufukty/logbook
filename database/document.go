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
			"created_at"`
	err := pool.QueryRow(
		context.Background(),
		query,
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
			&document.CreatedAt,
		)
	}
	return documents, err
}
