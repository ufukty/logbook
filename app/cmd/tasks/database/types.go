package database

import "regexp"

type (
	UserId      string
	ObjectiveId string 
	VersionId   string
	CommitId    string
	ActionId    string
	LinkId      string

	NonNegativeNumber int
)

var uuid = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func (v UserId) Validate() bool {
	return uuid.MatchString(string(v))
}

func (v ObjectiveId) Validate() bool {
	return uuid.MatchString(string(v))
}

func (v VersionId) Validate() bool {
	return uuid.MatchString(string(v))
}

func (v CommitId) Validate() bool {
	return uuid.MatchString(string(v))
}

func (v ActionId) Validate() bool {
	return uuid.MatchString(string(v))
}

func (v LinkId) Validate() bool {
	return uuid.MatchString(string(v))
}

func (v NonNegativeNumber) Validate() bool {
	return v >= 0
}

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

type LinkType string

const (
	Primary = LinkType("PRIMARY") // eg. When task owner break downs it
	Remote  = LinkType("REMOTE")  // eg. Collaborated objective attached to local objectives
	Private = LinkType("PRIVATE") //
)
