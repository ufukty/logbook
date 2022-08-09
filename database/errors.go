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
	ErrCreateUser = errors.New("ErrCreateUser")
)

var (
	ErrCreateDocument                 = errors.New("ErrCreateDocument")
	ErrGetDocumentByDocumentId        = errors.New("ErrGetDocumentByDocumentId")
	ErrGetHierarchicalViewItemsQuery  = errors.New("ErrGetHierarchicalViewItemsQuery")
	ErrGetHierarchicalViewItemsScan   = errors.New("ErrGetHierarchicalViewItemsScan")
	ErrGetChronologicalViewItemsQuery = errors.New("ErrGetChronologicalViewItemsQuery")
	ErrGetChronologicalViewItemsScan  = errors.New("ErrGetChronologicalViewItemsScan")
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

func CheckTaskId(taskId string) []error { // TODO:
	return nil
}
