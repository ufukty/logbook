package database

import "regexp"

type (
	Did               string // Document Id
	Iid               string // Item Id
	NonNegativeNumber int
)

var uuid = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func (v Did) Validate() bool {
	return uuid.MatchString(string(v))
}

func (v Iid) Validate() bool {
	return uuid.MatchString(string(v))
}

func (v NonNegativeNumber) Validate() bool {
	return v >= 0
}
