package database

import (
	"database/sql"
	"time"
)

type Document struct {
	DocumentId string    `json:"document_id"`
	UserId     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	ActiveTask string    `json:"active_task"`
}

type Bookmark struct {
	UserId       string         `gorm:""`
	TaskId       string         `gorm:""`
	DisplayName  string         `gorm:""`
	RootBookmark bool           `gorm:"default:FALSE"`
	CreatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    sql.NullString `gorm:""`
}
