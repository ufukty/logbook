package database

import (
	"time"
)

type TaskOperationSummary string

const (
	CREATE                    = TaskOperationSummary("CREATE")
	REORDER                   = TaskOperationSummary("REORDER")
	DELETE                    = TaskOperationSummary("DELETE")
	CONTENT_EDIT              = TaskOperationSummary("CONTENT_EDIT")
	NOTE_CREATE               = TaskOperationSummary("NOTE_CREATE")
	NOTE_EDIT                 = TaskOperationSummary("NOTE_EDIT")
	NOTE_DELETE               = TaskOperationSummary("NOTE_DELETE")
	MARK_COMPLETE             = TaskOperationSummary("MARK_COMPLETE")
	MARK_UNCOMPLETE           = TaskOperationSummary("MARK_UNCOMPLETE")
	COLLABORATION_ASSIGN      = TaskOperationSummary("COLLABORATION_ASSIGN")
	COLLABORATION_UNASSIGN    = TaskOperationSummary("COLLABORATION_UNASSIGN")
	COLLABORATION_RESTRICT    = TaskOperationSummary("COLLABORATION_RESTRICT")
	COLLABORATION_DERESTRICT  = TaskOperationSummary("COLLABORATION_DERESTRICT")
	COLLABORATION_CHANGE_ROLE = TaskOperationSummary("COLLABORATION_CHANGE_ROLE")
	HISTORY_ROLLBACK          = TaskOperationSummary("HISTORY_ROLLBACK")
	HISTORY_FASTFORWARD       = TaskOperationSummary("HISTORY_FASTFORWARD")
)

type TaskOperationStatus string

const (
	IN_REVIEW                   = TaskOperationStatus("IN_REVIEW")
	PRIV_ACCEPTED               = TaskOperationStatus("PRIV_ACCEPTED")
	PRIV_REJECTED               = TaskOperationStatus("PRIV_REJECTED")
	APPLIED_FASTFORWARD         = TaskOperationStatus("APPLIED_FASTFORWARD")
	CONFLICT_DETECTED           = TaskOperationStatus("CONFLICT_DETECTED")
	MANAGER_SELECTION_IN_REVIEW = TaskOperationStatus("MANAGER_SELECTION_IN_REVIEW")
	MANAGER_SELECTION_ACCEPTED  = TaskOperationStatus("MANAGER_SELECTION_ACCEPTED")
	MANAGER_SELECTION_APPLIED   = TaskOperationStatus("MANAGER_SELECTION_APPLIED")
	MANAGER_SELECTION_REJECTED  = TaskOperationStatus("MANAGER_SELECTION_REJECTED")
)

type User struct {
	UserId         string    `gorm:"<-:create;->;not null;unique;default:gen_random_UUID()"`
	NameSurname    string    `gorm:"<-:create;->;not null"`
	EmailAddress   string    `gorm:"<-:create;->;not null;unique"`
	Salt           string    `gorm:"<-:create;->;not null"`
	HashedPassword string    `gorm:"<-:create;->;not null"`
	Activated      bool      `gorm:"not null;default:FALSE"`
	ActivatedAt    time.Time `gorm:""`
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type TaskOperation struct {
	OperationId        string               `gorm:""`
	RevisionId         string               `gorm:""`
	PreviousRevisionId string               `gorm:""`
	TaskId             string               `gorm:""`
	Summary            TaskOperationSummary `gorm:""`
	Status             TaskOperationStatus  `gorm:""`
	UserId             string               `gorm:""`
	CreatedAt          time.Time            `gorm:"not null;default:CURRENT_TIMESTAMP"`
	ArchivedAt         time.Time            `gorm:""`
}

type Task struct {
	RevisionId            string    `gorm:"<-:create;->;not null;primaryKey"`
	TaskId                string    `gorm:"<-:create;->;not null;default:gen_random_UUID()"`
	OriginalCreatorUserId string    `gorm:"not null"`
	ResponsibleUserId     string    `gorm:"not null"`
	Content               string    `gorm:"not null"`
	Notes                 string    `gorm:""`
	CreatedAt             time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	CompletedAt           time.Time `gorm:""`
	ArchivedAt            time.Time `gorm:""`
	Archived              bool      `gorm:""`
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

type Bookmark struct {
	UserId       string    `gorm:""`
	TaskId       string    `gorm:""`
	DisplayName  string    `gorm:""`
	RootBookmark string    `gorm:""`
	CreatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    string    `gorm:""`
}

type TaskPermission struct {
}

func (Task) TableName() string {
	return `"TASK"`
}

func (User) TableName() string {
	return `"USER"`
}
