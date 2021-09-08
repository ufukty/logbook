package database

import (
	"time"
)

type Task struct {
	TaskId        string    `json:"task_id"`
	DocumentId    string    `json:"document_id"`
	ParentId      string    `json:"parent_id"`
	Content       string    `json:"content"`
	Degree        int       `json:"degree"`
	Depth         int       `json:"depth"`
	CreatedAt     time.Time `json:"created_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	ReadyToPickUp bool      `json:"ready_to_pick_up"`
}

type Document struct {
	DocumentId string    `json:"document_id"`
	CreatedAt  time.Time `json:"created_at"`
	ActiveTask *string   `json:"active_task"`
}
