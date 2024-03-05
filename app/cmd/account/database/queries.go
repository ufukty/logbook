package database

import (
	"context"
	"fmt"
)

func (db *Database) CreateUser() (*User, error) {
	user := &User{}
	query := `
	INSERT INTO "USER"
	DEFAULT VALUES
	RETURNING
		"uid", 
		"created_at", 
		COALESCE("active_task", '00000000-0000-0000-0000-000000000000')
`
	err := db.pool.QueryRow(context.Background(), query).Scan(
		&user.UserId,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning query result: %w", err)
	}
	return user, nil
}
