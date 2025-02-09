package columns

import (
	"logbook/internal/web/validate"
	"regexp"
)

func (v NonNegativeNumber) Validate() error {
	if v < 0 {
		return validate.ErrPattern
	}
	return nil
}

type pmm struct {
	pattern  *regexp.Regexp
	min, max int
}

func (b pmm) Validate(v string) error {
	return validate.StringBasics(v, b.min, b.max, b.pattern)
}

func r(s string) *regexp.Regexp {
	return regexp.MustCompile(s)
}

var (
	// creditCard   = pmm{r(`^(?:4[0-9]{12}(?:[0-9]{3})?)$`), 13, 19}
	// htmlDate     = pmm{r(`^\d{4}-\d{2}-\d{2}$`), 6, 8}
	// htmlDatetime = pmm{r(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$`), 9, 13}
	// htmlTime     = pmm{r(`^\d{2}:\d{2}$`), 3, 5}
	// numeric      = pmm{r(`^[1-9][0-9]*$`), 0, 100}
	// text         = pmm{r(`^[\p{L}0-9 ,.?!'’“”-]+$`), 0, 10000}
	// url          = pmm{r(`^[\p{L}0-9._%+-]+@[\p{L}0-9.-]+\.[\p{L}]{2,}$`), 0, 10000}
	email        = pmm{r(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`), 6, 150}
	groupTitle   = pmm{r(`[\p{L} ]+`), 2, 100}
	humanName    = pmm{r(`^\p{L}+([ '-]\p{L}+)*$`), 6, 100}
	phoneNumber  = pmm{r(`^\+?(\d{1,3})?[ -]?(\d{3})[ -]?(\d{3})[ -]?(\d{4})$`), 10, 15}
	sessionToken = pmm{r(`[A-Za-z0-9-_]+$`), 256, 256} // pattern is as defined in std lib base64.URLEncoding
	username     = pmm{r(`^[a-zA-Z]+[a-zA-Z0-9\_\.\-]*$`), 3, 50}
	uuid_        = pmm{r(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), len("00000000-0000-0000-0000-000000000000"), len("00000000-0000-0000-0000-000000000000")}
)

func (v AccessId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v BookmarkId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v BottomUpPropsId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v CollaborationId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v CollaboratorId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v ControlAreaId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v DelegationId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v Email) Validate() error {
	return email.Validate(string(v))
}

func (v GroupId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v GroupInviteId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v GroupMembershipId) Validate() error {
	return groupTitle.Validate(string(v))
}

func (v GroupName) Validate() error {
	return groupTitle.Validate(string(v))
}

func (v HumanName) Validate() error {
	return humanName.Validate(string(v))
}

func (v LinkId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v LoginId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v ObjectiveId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v OperationId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v Phone) Validate() error {
	return phoneNumber.Validate(string(v))
}

func (v PropertiesId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v SessionId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v SessionToken) Validate() error {
	return sessionToken.Validate(string(v))
}

func (v TagId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v UserId) Validate() error {
	return uuid_.Validate(string(v))
}

func (v Username) Validate() error {
	return username.Validate(string(v))
}

func (v VersionId) Validate() error {
	return uuid_.Validate(string(v))
}
