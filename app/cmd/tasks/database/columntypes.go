package database

import (
	"logbook/internal/web/validate"
	"regexp"

	"github.com/jackc/pgx/v5/pgtype"
)

type (
	UserId      string
	ObjectiveId string
	VersionId   string
	CommitId    string
	OperationId string
	LinkId      string
	HumanName   string

	NonNegativeNumber int
)

var (
	regexp_uuid         = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	regexp_text         = regexp.MustCompile(`^[\p{L}0-9 ,.?!'’“”-]+$`)
	regexp_human_name   = regexp.MustCompile(`^\p{L}+([ '-]\p{L}+)*$`)
	regexp_url          = regexp.MustCompile(`^[\p{L}0-9._%+-]+@[\p{L}0-9.-]+\.[\p{L}]{2,}$`)
	regexp_email        = regexp.MustCompile(`^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([\/\w .-]*)*\/?$`)
	regexp_username     = regexp.MustCompile(`^[A-Za-z0-9_]{3,15}$`)
	regexp_phone_number = regexp.MustCompile(`^\+?(\d{1,3})?[ -]?(\d{3})[ -]?(\d{3})[ -]?(\d{4})$`)
	regexp_date         = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`) // FIXME:
	regexp_numeric      = regexp.MustCompile(`^[1-9][0-9]*$`)
	regexp_credit_card  = regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?)$`)
)

var (
	max_length_uuid       = len("00000000-0000-0000-0000-000000000000")
	max_length_human_name = 100
)

var (
	min_length_uuid       = len("00000000-0000-0000-0000-000000000000")
	min_length_human_name = 1
)

func (v HumanName) Validate() error {
	return validate.StringBasics(string(v), min_length_human_name, max_length_human_name, regexp_human_name)
}

func (v UserId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v ObjectiveId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v VersionId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v CommitId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v OperationId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v LinkId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

func (v NonNegativeNumber) Validate() error {
	if v < 0 {
		return validate.ErrPattern
	}
	return nil
}

type LinkType string

const (
	Primary = LinkType("PRIMARY") // eg. When task owner break downs it
	Remote  = LinkType("REMOTE")  // eg. Collaborated objective attached to local objectives
	Private = LinkType("PRIVATE") //
)

const (
	ZeroObjectId  = ObjectiveId("00000000-0000-0000-0000-000000000000")
	ZeroVersionId = VersionId("00000000-0000-0000-0000-000000000000")
	ZeroUserId    = UserId("00000000-0000-0000-0000-000000000000")
	ZeroLinkId    = LinkId("00000000-0000-0000-0000-000000000000")
)

var ZeroTimestamp = pgtype.Timestamp{}
