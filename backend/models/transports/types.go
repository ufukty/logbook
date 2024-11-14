package transports

import (
	"time"
)

// account
type (
	HumanBirthday time.Time
	Password      string
	PhoneGrant    string
	EmailGrant    string
	PasswordGrant string
)

// all-encoded
type AntiCsrfToken string
