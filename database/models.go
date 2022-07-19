package database

import (
	"time"

	"github.com/jackc/pgtype"
)

type Task struct {
	TaskId string `json:"task_id"`

	UserId     string `json:"user_id"`
	DocumentId string `jsojn:"document_id"`

	Content string `json:"content"`

	ParentId      string `json:"parent_id"`
	Degree        int    `json:"degree"`
	Depth         int    `json:"depth"`
	Index         int    `json:"index"`
	ReadyToPickUp bool   `json:"ready_to_pick_up"`

	CompletionPct float64     `json:"completion_pct"`
	CompletedAt   pgtype.Date `json:"completed_at"` // nullable type

	CreatedAt time.Time `json:"created_at"`

	Archived bool `json:"archived"`
	Fold     bool `json:"fold"`
}

type Document struct {
	UserId     string    `json:"user_id"`
	DocumentId string    `json:"document_id"`
	CreatedAt  time.Time `json:"created_at"`
	ActiveTask string    `json:"active_task"`
}
