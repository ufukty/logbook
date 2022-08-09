package database

import "context"

func UserCreate(username string, emailAddress string, emailAddressTruncated string, passwordHashEncoded string) (*User, []error) {
	user := User{}
	query := `
		INSERT INTO "USER" (
			"username",
			"email_address",
			"email_address_truncated",
			"password_hash_encoded"
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			"user_id", 
			"username",
			"email_address",
			"email_address_truncated",
			"password_hash_encoded",
			"activated",
			"created_at"`
	err := pool.QueryRow(
		context.Background(),
		query,
		username,
		emailAddress,
		emailAddressTruncated,
		passwordHashEncoded,
	).Scan(
		&user.UserId,
		&user.Username,
		&user.EmailAddress,
		&user.EmailAddressTruncated,
		&user.PasswordHashEncoded,
		&user.Activated,
		&user.CreatedAt,
	)
	if err != nil {
		return &user, []error{err, ErrCreateUser}
	}
	return &user, nil
}

func UserGetByUserId(userId string) (*User, []error) {
	// TODO:
	user := User{}
	return &user, nil
}

func UserGetByUserName(userId string) (*User, []error) {
	// TODO:
	user := User{}
	return &user, nil
}
