package columns

import (
	"logbook/internal/web/validate"
	"regexp"
)

type (
	Username     string
	UserId       string
	LoginId      string
	AccessId     string
	SessionId    string
	SessionToken string
	Email        string
	HumanName    string
)

var (
	regexp_email      = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	regexp_human_name = regexp.MustCompile(`^\p{L}+([ '-]\p{L}+)*$`)
	regexp_username   = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9\_\.\-]*$`)
)

var (
	max_length_email         = 150
	max_length_human_name    = 100
	max_length_session_token = 256
	max_length_username      = 20
)

var (
	min_length_email         = 6
	min_length_human_name    = 6
	min_length_session_token = 256
	min_length_username      = 6
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

func (v SessionToken) Validate() error {
	return validate.StringBasics(string(v), min_length_session_token, max_length_session_token, regexp_base64_url)
}

func (v Email) Validate() error {
	return validate.StringBasics(string(v), min_length_email, max_length_email, regexp_email)
}

func (v HumanName) Validate() error {
	return validate.StringBasics(string(v), min_length_human_name, max_length_human_name, regexp_human_name)
}

const (
	ZeroUserId = UserId("00000000-0000-0000-0000-000000000000")
)

func (st *SessionToken) Set(v string) error {
	*st = SessionToken(v)
	return nil
}
