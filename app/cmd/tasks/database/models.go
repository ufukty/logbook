package database

import (
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
)

type OperationSummary string

const (
	TASK_CREATE               OperationSummary = "TASK_CREATE"
	TASK_REORDER              OperationSummary = "TASK_REORDER"
	TASK_DELETE               OperationSummary = "TASK_DELETE"
	TASK_CONTENT_EDIT         OperationSummary = "TASK_CONTENT_EDIT"
	TASK_MARK_COMPLETE        OperationSummary = "TASK_MARK_COMPLETE"
	TASK_MARK_UNCOMPLETE      OperationSummary = "TASK_MARK_UNCOMPLETE"
	NOTE_CREATE               OperationSummary = "NOTE_CREATE"
	NOTE_EDIT                 OperationSummary = "NOTE_EDIT"
	NOTE_DELETE               OperationSummary = "NOTE_DELETE"
	COLLABORATION_ASSIGN      OperationSummary = "COLLABORATION_ASSIGN"
	COLLABORATION_UNASSIGN    OperationSummary = "COLLABORATION_UNASSIGN"
	COLLABORATION_RESTRICT    OperationSummary = "COLLABORATION_RESTRICT"
	COLLABORATION_DERESTRICT  OperationSummary = "COLLABORATION_DERESTRICT"
	COLLABORATION_CHANGE_ROLE OperationSummary = "COLLABORATION_CHANGE_ROLE"
	HISTORY_ROLLBACK          OperationSummary = "HISTORY_ROLLBACK"
	HISTORY_FASTFORWARD       OperationSummary = "HISTORY_FASTFORWARD"
)

type OperationStatus string

const (
	SERVER_ORIGINATED           OperationStatus = "SERVER_ORIGINATED"
	IN_REVIEW                   OperationStatus = "IN_REVIEW"
	PRIV_ACCEPTED               OperationStatus = "PRIV_ACCEPTED"
	PRIV_REJECTED               OperationStatus = "PRIV_REJECTED"
	APPLIED_FASTFORWARD         OperationStatus = "APPLIED_FASTFORWARD"
	CONFLICT_DETECTED           OperationStatus = "CONFLICT_DETECTED"
	MANAGER_SELECTION_IN_REVIEW OperationStatus = "MANAGER_SELECTION_IN_REVIEW"
	MANAGER_SELECTION_ACCEPTED  OperationStatus = "MANAGER_SELECTION_ACCEPTED"
	MANAGER_SELECTION_APPLIED   OperationStatus = "MANAGER_SELECTION_APPLIED"
	MANAGER_SELECTION_REJECTED  OperationStatus = "MANAGER_SELECTION_REJECTED"
)

type Operation struct {
	OperationId         string           `gorm:"<-:create;->;unique;default:gen_random_UUID()"`
	PreviousOperationId sql.NullString   `gorm:"<-:create;->"`
	UserId              string           `gorm:"<-:create;->;not null"`
	Summary             OperationSummary `gorm:"<-:create;->;not null"`
	Status              OperationStatus  `gorm:"<-;->;not null"`
	CreatedAt           time.Time        `gorm:"<-:create;->;not null;default:CURRENT_TIMESTAMP"`
	ArchivedAt          time.Time        `gorm:"<-:->"`
}

type Task struct {
	TaskId        Iid               `json:"task_id"`
	DocumentId    Did               `json:"document_id"`
	ParentId      Iid               `json:"parent_id"`
	Content       string            `json:"content"`
	Degree        NonNegativeNumber `json:"degree"`
	Depth         NonNegativeNumber `json:"depth"`
	CreatedAt     time.Time         `json:"created_at"`
	CompletedAt   pgtype.Date       `json:"completed_at"` // nullable type
	ReadyToPickUp bool              `json:"ready_to_pick_up"`

	RevisionId            string    `gorm:"<-:create;->;not null;primaryKey"`
	TaskId                string    `gorm:"<-:create;->;not null;default:gen_random_UUID()"`
	OriginalCreatorUserId string    `gorm:"not null"`
	ResponsibleUserId     string    `gorm:"not null"`
	Content               string    `gorm:"not null"`
	CreatedAt             time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	CompletedAt           time.Time `gorm:""`
	ArchivedAt            time.Time `gorm:""`
	Archived              bool      `gorm:""`
	RootTask              bool      `gorm:"default:FALSE"`
}

type TaskLink struct {
	LinkId              string    `gorm:""`
	RevisionId          string    `gorm:""`
	TaskId              string    `gorm:""`
	TaskRevisionId      string    `gorm:""`
	SuperTaskId         string    `gorm:""`
	SuperTaskRevisionId string    `gorm:""`
	Index               string    `gorm:""`
	PrimaryLink         string    `gorm:""`
	CreatedAt           time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type TaskLinkUserPreferences struct {
	LinkId string `gorm:""`
	UserId string `gorm:""`
	Fold   bool   `gorm:""`
}

type TaskProps struct {
	RevisionId           string  `gorm:""`
	TaskId               string  `gorm:""`
	UserId               string  `gorm:""`
	Degree               int     `gorm:""`
	Depth                int     `gorm:""`
	CompletionPercentage float64 `gorm:""`
}

type TaskPermission struct {
}
