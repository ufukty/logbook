package transports

import "logbook/internal/web/validate"

func (v Password) Validate() error {
	return validate.StringBasics(string(v), min_length_password, max_length_password, nil)
}
