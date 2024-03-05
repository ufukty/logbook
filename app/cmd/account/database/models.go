package database

import "time"

type User struct {
	UserId            string
	NameSurname       string
	EmailAddress      string
	SaltBase64Encoded string
	HashEncoded       string
	Activated         bool
	ActivatedAt       time.Time
	CreatedAt         time.Time
}
