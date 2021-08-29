package database

import (
	"time"
)

type TaskStatus string

const (
	Archive      TaskStatus = "archive"
	Active       TaskStatus = "active"
	Paused       TaskStatus = "paused"
	ReadyToStart TaskStatus = "ready-to-start"
	Drawer       TaskStatus = "drawer"
)

func StringToTaskStatus(str string) (TaskStatus, error) {
	switch str {
	case "archive":
		return Archive, nil
	case "active":
		return Active, nil
	case "paused":
		return Paused, nil
	case "ready-to-start":
		return ReadyToStart, nil
	case "drawer":
		return Drawer, nil
	}
	return "", ErrStrToTaskStatusNoMatchingRecord
}

type Task struct {
	Content     string     `json:"content"`
	CreatedAt   time.Time  `json:"created_at"`
	Degree      int        `json:"degree"`
	Depth       int        `json:"depth"`
	ParentId    string     `json:"parent_id"`
	TaskGroupId string     `json:"task_group_id"`
	TaskId      string     `json:"task_id"`
	TaskStatus  TaskStatus `json:"task_status"`
}

type TaskGroup struct {
	CreatedAt     time.Time  `json:"created_at"`
	DocumentId    string     `json:"document_id"`
	TaskGroupId   string     `json:"task_group_id"`
	TaskGroupType TaskStatus `json:"task_group_type"`
	Tasks         []Task     `json:"tasks"`
	TotalTasks    int        `json:"total_tasks"`
}

type Document struct {
	TaskGroups      []TaskGroup `json:"task_groups"`
	TotalTaskGroups int         `json:"total_task_groups"`
	CreatedAt       time.Time   `json:"created_at"`
	DocumentId      string      `json:"document_id"`
}

type DocumentReference struct { // FIXME: REMOVE
	DocumentID  string    `json:"document_id"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}
