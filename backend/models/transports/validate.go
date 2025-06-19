package transports

import (
	"time"

	"github.com/ufukty/gohandlers/pkg/validator/validate"
)

func (a AntiCsrfToken) Validate() any {
	return validate.String(string(a), length_anti_csrf_token, length_anti_csrf_token, regexp_base64_url)
}

func (hb HumanBirthday) Validate() any {
	return validate.Time(time.Time(hb), min_human_birthday, max_human_birthday)
}

func (v Password) Validate() any {
	return validate.String(string(v), min_length_password, max_length_password, nil)
}

func (v PhoneGrant) Validate() any {
	return validate.String(string(v), length_uuid, length_uuid, regexp_uuid)
}

func (v EmailGrant) Validate() any {
	return validate.String(string(v), length_uuid, length_uuid, regexp_uuid)
}

func (v PasswordGrant) Validate() any {
	return validate.String(string(v), length_uuid, length_uuid, regexp_uuid)
}
