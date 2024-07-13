package requests

import (
	"fmt"
	"strconv"
)

type UserUuid string

func (id *UserUuid) Set(v string) error {
	*id = UserUuid(v)
	return nil
}

type UserId int

func (id *UserId) Set(v string) error {
	i, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("processing text for the number value: %w", err)
	}
	*id = UserId(i)
	return nil
}

type MockStringType string

func (m *MockStringType) Set(s string) error {
	*m = MockStringType(s)
	return nil
}

type MockIntType int

func (m *MockIntType) Set(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*m = MockIntType(val)
	return nil
}

type TestRequest struct {
	SessionToken Cookie[MockIntType] `cookie:"session_token"`
	UserID       MockIntType         `url:"user_id"`
	Name         MockStringType      `json:"name"`
}
