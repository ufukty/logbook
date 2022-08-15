package parameters

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type TypeCheckable interface {
	TypeCheck() error
}

func (s *NonEmptyString) TypeCheck() error {
	if (*s) == "" {
		return errors.New("NonEmptyString is empty")
	} else {
		return nil
	}
}

func (s *UserId) TypeCheck() error {
	if (*s) == "" {
		return errors.New("UserId is empty")
	}
	if _, err := uuid.Parse((string)(*s)); err != nil {
		return errors.Wrap(err, "UserId is not valid uuid")
	}
	return nil
}

func (s *TaskId) TypeCheck() error {
	if (*s) == "" {
		return errors.New("TaskId is empty")
	}
	if _, err := uuid.Parse((string)(*s)); err != nil {
		return errors.Wrap(err, "TaskId is not valid uuid")
	}
	return nil
}
