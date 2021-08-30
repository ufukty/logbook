package database

import "errors"

var ErrNoResult error

var ErrInvalidInput error

var ErrEmptyDocumentId error

var ErrInvalidDocumentId error

var ErrEmptyTaskGroupId error

var ErrInvalidTaskGroupId error

var ErrStrToTaskStatusNoMatchingRecord error

var (
	ErrCreateDocument                   = errors.New("CreateDocument faced with an error")
	ErrCreateDocumentWithTaskGroups     = errors.New("CreateDocumentWithTaskGroups faced with an error")
	ErrGetDocumentByDocumentId          = errors.New("GetDocumentByDocumentId faced with an error")
	ErrCreateTaskGroup                  = errors.New("CreateTaskGroup faced with an error")
	ErrGetTaskGroupByTaskGroupId        = errors.New("GetTaskGroupByTaskGroupId faced with an error")
	ErrGetTaskGroupsByDocumentIdQuery   = errors.New("GetTaskGroupsByDocumentIdQuery faced with an error")
	ErrGetTaskGroupsByDocumentIdRowScan = errors.New("GetTaskGroupsByDocumentIdRowScan faced with an error")
	ErrCreateTask                       = errors.New("CreateTask faced with an error")
	ErrGetTaskByTaskId                  = errors.New("GetTaskByTaskId faced with an error")
	ErrGetTasksByTaskGroupIdQuery       = errors.New("GetTasksByTaskGroupIdQuery faced with an error")
	ErrGetTasksByTaskGroupIdScan        = errors.New("GetTasksByTaskGroupIdScan faced with an error")
	ErrGetTaskByParentIdQuery           = errors.New("GetTaskByParentIdQuery faced with an error")
	ErrGetTaskByParentIdScan            = errors.New("GetTaskByParentIdScan faced with an error")
	ErrUpdateTaskItem                   = errors.New("UpdateTaskItem faced with an error")
)

func CheckDocumentId(documentId string) []error {
	// check existance
	if documentId == "" {
		return []error{ErrEmptyDocumentId}
	}
	// check validity
	if _, err := GetDocumentByDocumentId(documentId); err != nil {
		return append(err, ErrInvalidDocumentId)
	}
	return nil
}

func CheckTaskGroupId(taskGroupId string) []error {
	// check existance
	if taskGroupId == "" {
		return []error{ErrEmptyTaskGroupId}
	}
	// check validity
	if _, err := GetTaskGroupByTaskGroupId(taskGroupId); err != nil {
		return append(err, ErrInvalidTaskGroupId)
	}
	return nil
}

func CheckTaskId(taskId string) []error { // TODO:
	return nil
}
