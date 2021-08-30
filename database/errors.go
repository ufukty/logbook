package database

import "errors"

var ErrNotSpecified error

var ErrNoResult error

var ErrInvalidInput error

var ErrEmptyDocumentId error

var ErrInvalidDocumentId error

var ErrEmptyTaskGroupId error

var ErrInvalidTaskGroupId error

var ErrStrToTaskStatusNoMatchingRecord = errors.New("")

// FIXME: Wrap the PSQL error with Err class

func initErrors() {
	ErrNotSpecified = errors.New("ErrNotSpecified")
	ErrNoResult = errors.New("ErrNoResult")
	ErrInvalidInput = errors.New("ErrInvalidInput")
}

func checkDocumentId(documentId string) error {
	// check existance
	if documentId == "" {
		return ErrEmptyDocumentId
	}
	// check validity
	if _, err := GetDocumentByDocumentId(documentId); err != nil {
		return err
	}
	return nil
}

func checkTaskGroupId(taskGroupId string) error {
	// check existance
	if taskGroupId == "" {
		return ErrEmptyTaskGroupId
	}
	// check validity
	if _, err := GetTaskGroupByTaskGroupId(taskGroupId); err != nil {
		return err
	}
	return nil
}

func checkTaskId(taskId string) error { // TODO:
	return nil
}

func exportError(err error) error {
	if err == nil {
		return nil
	}
	switch err.Error() {
	case "no rows in result set":
		return ErrNoResult
	default:
		return ErrNotSpecified
	}
}
