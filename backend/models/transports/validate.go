package transports

import (
	"logbook/models/validators"
	"time"

	"github.com/ufukty/gohandlers/pkg/validator/validate"
)

func (v EmailGrant) Validate() any    { return validators.Uuid.Validate(string(v)) }
func (v PasswordGrant) Validate() any { return validators.Uuid.Validate(string(v)) }
func (v PhoneGrant) Validate() any    { return validators.Uuid.Validate(string(v)) }

func (v AntiCsrfToken) Validate() any { return validators.AntiCsrfToken.Validate(string(v)) }
func (v HumanBirthday) Validate() any { return validate.Time(time.Time(v), minBirthday, maxBirthday) }
func (v Password) Validate() any      { return validate.String(string(v), 6, 2048, nil) }

var (
	minBirthday = time.Now().AddDate(-100, 0, 0) // 100 years ago
	maxBirthday = time.Now().AddDate(-18, 0, 1)  // 18 years ago
)
