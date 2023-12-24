package database

import "regexp"

type (
	Username          string
	Uid               string // User Id
	NonNegativeNumber int
)

var uuid = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
var username = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9\_\.\-]*$`)

func (v Username) Validate() bool {
	return username.MatchString(string(v))
}

func (v Uid) Validate() bool {
	return uuid.MatchString(string(v))
}

func (v NonNegativeNumber) Validate() bool {
	return v >= 0
}
