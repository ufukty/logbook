package columns

import (
	"logbook/internal/web/validate"
	"regexp"
)

type (
	ObjectiveId string
	VersionId   string
	CommitId    string
	OperationId string
	LinkId      string
)

var (
	regexp_phone_number = regexp.MustCompile(`^\+?(\d{1,3})?[ -]?(\d{3})[ -]?(\d{3})[ -]?(\d{4})$`)
	regexp_credit_card  = regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?)$`)
)

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

type LinkType string

const (
	Primary = LinkType("PRIMARY") // eg. When task owner break downs it
	Remote  = LinkType("REMOTE")  // eg. Collaborated objective attached to local objectives
	Private = LinkType("PRIVATE") //
)

const (
	ZeroObjectId  = ObjectiveId("00000000-0000-0000-0000-000000000000")
	ZeroVersionId = VersionId("00000000-0000-0000-0000-000000000000")
	ZeroLinkId    = LinkId("00000000-0000-0000-0000-000000000000")
)
