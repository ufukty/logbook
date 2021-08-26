package database

import "context"

func createUser(user User) (User, bool, error) {
	query := `
		INSERT INTO "USER" (
			"email",
			"password"
		) 
		VALUES (
			$1, $2
		)
		RETURNING
			"user_id",
			"created_at"`
	err := pool.QueryRow(
		context.Background(),
		query,
		user.Email,
		user.Password,
	).Scan(
		&user.UserID,
		&user.CreatedAt,
	)
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

func getUserByUserId(userId string) (User, bool, error) {
	user := User{UserID: userId}
	query := `
		SELECT
			"email",
			"password",
			"created_at"
		FROM 
			"USER"
		WHERE
			"user_id"=$1`

	err := pool.QueryRow(
		context.Background(),
		query,
		userId,
	).Scan(
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}
