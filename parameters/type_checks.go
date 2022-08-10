package parameters

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
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

// func (s *EmailAddress) TypeCheck() error {
// 	if (*s) == "" {
// 		return errors.New("email doesn't pass type checking")
// 	} else {
// 		return nil
// 	}
// }

func (s *UserId) TypeCheck() error {
	if (*s) == "" {
		return errors.New("UserId is empty")
	}
	if _, err := uuid.Parse((string)(*s)); err != nil {
		return fmt.Errorf("UserId is not valid uuid \"%w\"", err)
	}
	return nil
}

func (s *TaskId) TypeCheck() error {
	if (*s) == "" {
		return errors.New("TaskId is empty")
	}
	if _, err := uuid.Parse((string)(*s)); err != nil {
		return fmt.Errorf("TaskId is not valid uuid: %w", err)
	}
	return nil
}
