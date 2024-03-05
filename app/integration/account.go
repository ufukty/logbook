package integration

import (
	"fmt"
	"logbook/cmd/account/endpoints"
)

func (ctl UserClient) Register() error {
	_, err := ctl.accctl.Register(&endpoints.CreateUserRequest{
		Email:       "test@localhost.xy",
		NameSurname: "Tést McSingleton",
		Password:    "123456789",
		Username:    "mcsingleton",
	})
	if err != nil {
		return fmt.Errorf("making the request: %w", err)
	}
	return nil
}
