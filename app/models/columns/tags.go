package columns

import (
	"logbook/internal/web/validate"
)

type (
	TagId string
)

func (v TagId) Validate() error {
	return validate.StringBasics(string(v), min_length_uuid, max_length_uuid, regexp_uuid)
}
