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

var (
	regexp_uuid         = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	regexp_text         = regexp.MustCompile(`^[\p{L}0-9 ,.?!'’“”-]+$`)
	regexp_human_name   = regexp.MustCompile(`^\p{L}+([ '-]\p{L}+)*$`)
	regexp_url          = regexp.MustCompile(`^[\p{L}0-9._%+-]+@[\p{L}0-9.-]+\.[\p{L}]{2,}$`)
	regexp_email        = regexp.MustCompile(`^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([\/\w .-]*)*\/?$`)
	regexp_username     = regexp.MustCompile(`^[A-Za-z0-9_]{3,15}$`)
	regexp_phone_number = regexp.MustCompile(`^\+?(\d{1,3})?[ -]?(\d{3})[ -]?(\d{3})[ -]?(\d{4})$`)
	regexp_date         = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`) // FIXME:
	regexp_numeric      = regexp.MustCompile(`^[1-9][0-9]*$`)
	regexp_credit_card  = regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?)$`)
)

func (v UserId) Validate() bool {
	return regexp_uuid.MatchString(string(v))
}

func (v ObjectiveId) Validate() bool {
	return regexp_uuid.MatchString(string(v))
}

func (v VersionId) Validate() bool {
	return regexp_uuid.MatchString(string(v))
}

func (v CommitId) Validate() bool {
	return regexp_uuid.MatchString(string(v))
}

func (v ActionId) Validate() bool {
	return regexp_uuid.MatchString(string(v))
}

func (v LinkId) Validate() bool {
	return regexp_uuid.MatchString(string(v))
}

func (v NonNegativeNumber) Validate() bool {
	return v >= 0
}

type ActionSummary string

const (
	ObjectiveCreate           ActionSummary = "OBJECTIVE_CREATE"
	ObjectiveReorder          ActionSummary = "OBJECTIVE_REORDER"
	ObjectiveDelete           ActionSummary = "OBJECTIVE_DELETE"
	ObjectiveTextAssign       ActionSummary = "TEXT_ASSIGN"
	ObjectiveMarkComplete     ActionSummary = "TASK_MARK_COMPLETE"
	ObjectiveUnmarkComplete   ActionSummary = "TASK_MARK_UNCOMPLETE"
	ObjectiveNoteAssign       ActionSummary = "OBJECTIVE_"
	CollaborationAssign       ActionSummary = "COLLABORATION_ASSIGN"
	CollaborationUnassign     ActionSummary = "COLLABORATION_UNASSIGN"
	COLLABORATION_RESTRICT    ActionSummary = "COLLABORATION_RESTRICT"
	COLLABORATION_DERESTRICT  ActionSummary = "COLLABORATION_DERESTRICT"
	COLLABORATION_CHANGE_ROLE ActionSummary = "COLLABORATION_CHANGE_ROLE"
	HISTORY_ROLLBACK          ActionSummary = "HISTORY_ROLLBACK"
	HISTORY_FASTFORWARD       ActionSummary = "HISTORY_FASTFORWARD"
)

type ActionStatus string

const (
	SERVER_ORIGINATED           ActionStatus = "SERVER_ORIGINATED"
	IN_REVIEW                   ActionStatus = "IN_REVIEW"
	PRIV_ACCEPTED               ActionStatus = "PRIV_ACCEPTED"
	PRIV_REJECTED               ActionStatus = "PRIV_REJECTED"
	APPLIED_FASTFORWARD         ActionStatus = "APPLIED_FASTFORWARD"
	CONFLICT_DETECTED           ActionStatus = "CONFLICT_DETECTED"
	MANAGER_SELECTION_IN_REVIEW ActionStatus = "MANAGER_SELECTION_IN_REVIEW"
	MANAGER_SELECTION_ACCEPTED  ActionStatus = "MANAGER_SELECTION_ACCEPTED"
	MANAGER_SELECTION_APPLIED   ActionStatus = "MANAGER_SELECTION_APPLIED"
	MANAGER_SELECTION_REJECTED  ActionStatus = "MANAGER_SELECTION_REJECTED"
)

type LinkType string

const (
	Primary = LinkType("PRIMARY") // eg. When task owner break downs it
	Remote  = LinkType("REMOTE")  // eg. Collaborated objective attached to local objectives
	Private = LinkType("PRIVATE") //
)

var NullObjectId = ObjectiveId("00000000-0000-0000-0000-000000000000")
