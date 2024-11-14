package app

import (
	"context"
	"fmt"

	profiles "logbook/cmd/profiles/endpoints"
	sessions "logbook/cmd/sessions/endpoints"
	"logbook/cmd/users/endpoints"
	"logbook/models/columns"
	"logbook/models/transports"
)

type RegisterRequest struct {
	AntiCsrfToken transports.AntiCsrfToken

	Firstname columns.HumanName
	Lastname  columns.HumanName
	Birthday  transports.HumanBirthday
	Country   transports.Country

	EmailGrant    transports.EmailGrant
	PhoneGrant    transports.PhoneGrant
	PasswordGrant transports.PasswordGrant
}

var ErrEmailExists = fmt.Errorf("email in use")

/*
 * Objectives for this function
 * DONE: Sanitize user input
 * DONE: Produce unique salt and hash user password with it
 * DONE: Secure against timing-attacks
 * TODO: Check anti-CSRF token
 * DONE: Check account duplication (attempt to register with same e-mail twice)
 * TODO: Create first task
 * TODO: Create privilege table record for created task
 * TODO: Create operation table record for task creation
 * TODO: Create first bookmark
 * TODO: Wrap creation of user-task-bookmark with transaction, rollback on failure to not-lock person to re-register with same email
 */
func (a *App) Register(ctx context.Context, params RegisterRequest) error {
	cu, err := a.Users.CreateUser(&endpoints.CreateUserRequest{})
	if err != nil {
		return fmt.Errorf("User.CreateUser: %w", err)
	}

	email, ok := a.grants.email.Get(params.EmailGrant)
	if !ok {
		return fmt.Errorf("invalid email grant")
	}

	password, ok := a.grants.password.Get(params.PasswordGrant)
	if !ok {
		return fmt.Errorf("invalid password grant")
	}

	_, err = a.Sessions.SaveCredentials(&sessions.SaveCredentialsRequest{
		Uid:      cu.Uid,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("a.Sessions.SaveCredentials: %w", err)
	}

	_, err = a.Profiles.CreateProfile(&profiles.CreateProfileRequest{
		Uid:       cu.Uid,
		Firstname: params.Firstname,
		Lastname:  params.Lastname,
	})
	if err != nil {
		return fmt.Errorf("profiles.CreateProfile: %w", err)
	}

	// _, err = q.SelectLatestLoginByEmail(ctx, params.Email)
	// if err == nil {
	// 	return ErrEmailExists
	// }

	// user, err := q.InsertUser(ctx)
	// if err != nil {
	// 	return fmt.Errorf("inserting record to database: %w", err)
	// }

	// _, err = q.InsertProfileInformation(ctx, database.InsertProfileInformationParams{
	// 	Uid:       user.Uid,
	// 	Firstname: params.Firstname,
	// 	Lastname:  params.Lastname,
	// })
	// if err != nil {
	// 	return fmt.Errorf("inserting profile information into database: %w", err)
	// }

	// _, err = a.Objectives.RockCreate(&endpoints.RockCreateRequest{
	// 	UserId: user.Uid,
	// })
	// if err != nil {
	// 	return fmt.Errorf("creating rock for user via objectives service: %w", err)
	// }

	// err = tx.Commit(ctx)
	// if err != nil {
	// 	return fmt.Errorf("tx.Commit: %w", err)
	// }
	return nil
}
