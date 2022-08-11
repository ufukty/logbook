package parameters

import "time"

type NonEmptyString string
type EmailAddress string
type UserId string
type TaskId string

type OperationType NonEmptyString
type OperationDetails NonEmptyString

//  keep below if you gonna include them in Parameter.Response; otherwise remove

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
	UserId                string    `json:"user_id"`
	Username              string    `json:"username"`
	EmailAddress          string    `json:"email_address"`
	EmailAddressTruncated string    `json:"email_address_truncated"`
	PasswordHashEncoded   string    `json:"password_hash_encoded"`
	Activated             bool      `json:"activated"`
	CreatedAt             time.Time `json:"created_at"`
}

type TaskOperation struct {
	OperationId        string               `json:"operation_id"`
	RevisionId         string               `json:"revision_id"`
	PreviousRevisionId string               `json:"previous_revision_id"`
	TaskId             string               `json:"task_id"`
	Summary            TaskOperationSummary `json:"operation_summary"`
	Status             TaskOperationStatus  `json:"operation_status"`
	UserId             string               `json:"user_id"`
	CreatedAt          time.Time            `json:"created_at"`
	ArchivedAt         time.Time            `json:"archived_at"`
}

type Task struct {
	RevisionId            string    `json:"revision_id"`
	TaskId                string    `json:"task_id"`
	OriginalCreatorUserId string    `json:"original_creator_user_id"`
	ResponsibleUserId     string    `json:"responsible_user_id"`
	Content               string    `json:"content"`
	Notes                 string    `json:"notes"`
	CreatedAt             time.Time `json:"created_at"`
	CompletedAt           time.Time `json:"completed_at"`
	ArchivedAt            time.Time `json:"archived_at"`
	Archived              bool      `json:"archived"`
}

type TaskLink struct {
	LinkId              string    `json:"link_id"`
	RevisionId          string    `json:"revision_id"`
	TaskId              string    `json:"task_id"`
	TaskRevisionId      string    `json:"task_revision_id"`
	SuperTaskId         string    `json:"super_task_id"`
	SuperTaskRevisionId string    `json:"super_task_revision_id"`
	Index               string    `json:"index"`
	PrimaryLink         string    `json:"primary_link"`
	CreatedAt           time.Time `json:"created_at"`
}

type TaskLinkUserPreferences struct {
	LinkId string `json:"link_id"`
	UserId string `json:"user_id"`
	Fold   bool   `json:"fold"`
}

type TaskProps struct {
	RevisionId           string  `json:"revision_id"`
	TaskId               string  `json:"task_id"`
	UserId               string  `json:"user_id"`
	Degree               int     `json:"degree"`
	Depth                int     `json:"depth"`
	CompletionPercentage float64 `json:"completion_percentage"`
}

type Bookmark struct {
	UserId       string    `json:"user_id"`
	TaskId       string    `json:"task_id"`
	DisplayName  string    `json:"display_name"`
	RootBookmark string    `json:"root_bookmark"`
	CreatedAt    time.Time `json:"created_at"`
	DeletedAt    string    `json:"deleted_at"`
}

type TaskPermission struct {
}
