package database

import "errors"

var ErrNoResult error

var ErrInvalidInput error

var ErrEmptyDocumentId error

var ErrInvalidDocumentId error

var ErrEmptyTaskGroupId error

var ErrInvalidTaskGroupId error

var ErrStrToTaskStatusNoMatchingRecord = errors.New("ErrStrToTaskStatusNoMatchingRecord")

var (
	ErrCreateDocument                         = errors.New("ErrCreateDocument")
	ErrGetDocumentByDocumentId                = errors.New("ErrGetDocumentByDocumentId")
	ErrGetDocumentOverviewWithDocumentIdQuery = errors.New("ErrGetDocumentOverviewWithDocumentIdQuery")
	ErrGetDocumentOverviewWithDocumentIdScan  = errors.New("ErrGetDocumentOverviewWithDocumentIdScan")
	ErrGetChronologicalViewItemsQuery         = errors.New("ErrGetChronologicalViewItemsQuery")
	ErrGetChronologicalViewItemsScan          = errors.New("ErrGetChronologicalViewItemsScan")
)

var (
	ErrCreateTask                = errors.New("ErrCreateTask")
	ErrGetTaskByTaskId           = errors.New("ErrGetTaskByTaskId")
	ErrGetTasksByDocumentIdQuery = errors.New("ErrGetTasksByDocumentIdQuery")
	ErrGetTasksByDocumentIdScan  = errors.New("ErrGetTasksByDocumentIdScan")
	ErrGetTaskByParentIdQuery    = errors.New("ErrGetTaskByParentIdQuery")
	ErrGetTaskByParentIdScan     = errors.New("ErrGetTaskByParentIdScan")
	ErrUpdateTaskItem            = errors.New("ErrUpdateTaskItem")
	ErrMarkATaskDone             = errors.New("ErrMarkATaskDone")
)

func (db *Database) CheckDocumentId(documentId string) []error {
	// check existance
	if documentId == "" {
		return []error{ErrEmptyDocumentId}
	}
	// check validity
	if _, err := db.GetDocumentByDocumentId(documentId); err != nil {
		return append(err, ErrInvalidDocumentId)
	}
	return nil
}

func CheckTaskId(taskId string) []error { // TODO:
	return nil
}
