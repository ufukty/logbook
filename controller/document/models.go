package document

import "time"

type Task struct {
	TaskId     string    `json:"task_id"`
	ParentId   string    `json:"parent_id"`
	Content    string    `json:"content"`
	TaskStatus string    `json:"task_status"`
	Degree     int       `json:"degree"`
	Depth      int       `json:"depth"`
	CreatedAt  time.Time `json:"created_at"`
}

type TaskGroupType string

const (
	Archive      TaskGroupType = "archive"
	Active       TaskGroupType = "active"
	Paused       TaskGroupType = "paused"
	ReadyToStart TaskGroupType = "ready-to-start"
	Drawer       TaskGroupType = "drawer"
)

type TaskGroup struct {
	GroupId    string        `json:"group_id"`
	GroupType  TaskGroupType `json:"group_type"`
	Tasks      []Task        `json:"tasks"`
	TotalTasks int           `json:"total_tasks"`
}

type Document struct {
	DocumentId      string      `json:"document_id"`
	TaskGroups      []TaskGroup `json:"task_groups"`
	TotalTaskGroups int         `json:"total_task_groups"`
}

type DocumentReference struct {
	DocumentID  string    `json:"document_id"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}

type Dashboard struct {
	UserId    string              `json:"user_id"`
	Documents []DocumentReference `json:"documents"`
}
