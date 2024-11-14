package transports

import (
	"fmt"
	"logbook/internal/web/validate"
	"time"
)

func (a AntiCsrfToken) Validate() error {
	return validate.StringBasics(string(a), length_anti_csrf_token, length_anti_csrf_token, regexp_base64_url)
}

func (hb HumanBirthday) Validate() error {
	if !validate.TimeBasics(time.Time(hb), min_human_birthday, max_human_birthday) {
		return fmt.Errorf("out of range")
	}
	return nil
}

func (v Password) Validate() error {
	return validate.StringBasics(string(v), min_length_password, max_length_password, nil)
}

func (v PhoneGrant) Validate() error {
	return validate.StringBasics(string(v), length_uuid, length_uuid, regexp_uuid)
}

func (v EmailGrant) Validate() error {
	return validate.StringBasics(string(v), length_uuid, length_uuid, regexp_uuid)
}

func (v PasswordGrant) Validate() error {
	return validate.StringBasics(string(v), length_uuid, length_uuid, regexp_uuid)
}
