package database

import (
	"time"

	"github.com/jackc/pgtype"
)

type Task struct {
	TaskId        string      `json:"task_id"`
	DocumentId    string      `json:"document_id"`
	ParentId      string      `json:"parent_id"`
	Content       string      `json:"content"`
	Degree        int         `json:"degree"`
	Depth         int         `json:"depth"`
	CreatedAt     time.Time   `json:"created_at"`
	CompletedAt   pgtype.Date `json:"completed_at"` // nullable type
	ReadyToPickUp bool        `json:"ready_to_pick_up"`
}

type Document struct {
	DocumentId string    `json:"document_id"`
	CreatedAt  time.Time `json:"created_at"`
	ActiveTask string    `json:"active_task"`
}
