package database

import (
	"logbook/internal/web/validate"
	"regexp"

	"github.com/jackc/pgx/v5/pgtype"
)

type (
	Username          string
	UserId            string
	LoginId           string
	AccessId          string
	SessionId         string
	Email             string
	NonNegativeNumber int
	HumanName         string
)

var (
	regexp_email      = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	regexp_human_name = regexp.MustCompile(`^\p{L}+([ '-]\p{L}+)*$`)
	regexp_username   = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9\_\.\-]*$`)
	regexp_uuid       = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

var (
	max_length_email      = 150
	max_length_human_name = 100
	max_length_username   = 20
	max_length_uuid       = len("00000000-0000-0000-0000-000000000000")
)

var (
	min_length_email      = 6
	min_length_human_name = 6
	min_length_username   = 6
	min_length_uuid       = len("00000000-0000-0000-0000-000000000000")
)

func (v Username) Validate() error {
	return validate.StringBasics(string(v), min_length_username, max_length_username, regexp_username)
}

func (v UserId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v LoginId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v AccessId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v SessionId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v Email) Validate() error {
	return validate.StringBasics(string(v), min_length_email, max_length_email, regexp_email)
}

func (v NonNegativeNumber) Validate() error {
	if v < 0 {
		return validate.ErrPattern
	}
	return nil
}

func (v HumanName) Validate() error {
	return validate.StringBasics(string(v), min_length_human_name, max_length_human_name, regexp_human_name)
}

const (
	ZeroUserId = UserId("00000000-0000-0000-0000-000000000000")
)

var ZeroTimestamp = pgtype.Timestamp{}
