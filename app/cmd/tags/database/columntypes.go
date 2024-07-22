package database

import (
	"logbook/internal/web/validate"
	"regexp"
)

type (
	TagId string
)

var (
	regexp_uuid = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

var (
	max_length_uuid = len("00000000-0000-0000-0000-000000000000")
)

var (
	min_length_uuid = len("00000000-0000-0000-0000-000000000000")
)

func (v TagId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}

const (
	ZeroUserId = TagId("00000000-0000-0000-0000-000000000000")
)
