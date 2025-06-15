package columns

import (
	"github.com/ufukty/gohandlers/pkg/validator"
	"github.com/ufukty/gohandlers/pkg/validator/validate"
)

func (v NonNegativeNumber) Validate() error {
	if v < 0 {
		return validate.ErrPattern
	}
	return nil
}

var (
	// creditCard   = validator.ForStrings(`^(?:4[0-9]{12}(?:[0-9]{3})?)$`, 13, 19)
	// htmlDate     = validator.ForStrings(`^\d{4}-\d{2}-\d{2}$`, 6, 8)
	// htmlDatetime = validator.ForStrings(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$`, 9, 13)
	// htmlTime     = validator.ForStrings(`^\d{2}:\d{2}$`, 3, 5)
	// numeric      = validator.ForStrings(`^[1-9][0-9]*$`, 0, 100)
	// text         = validator.ForStrings(`^[\p{L}0-9 ,.?!'’“”-]+$`, 0, 10000)
	// url          = validator.ForStrings(`^[\p{L}0-9._%+-]+@[\p{L}0-9.-]+\.[\p{L}]{2,}$`, 0, 10000)
	email        = validator.ForStrings(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, 6, 150)
	groupName    = validator.ForStrings(`[\p{L} ]+`, 2, 100)
	humanName    = validator.ForStrings(`^\p{L}+([ '-]\p{L}+)*$`, 6, 100)
	phoneNumber  = validator.ForStrings(`^\+?(\d{1,3})?[ -]?(\d{3})[ -]?(\d{3})[ -]?(\d{4})$`, 10, 15)
	sessionToken = validator.ForStrings(`[A-Za-z0-9-_]+$`, 256, 256) // pattern is as defined in std lib base64.URLEncoding
	username     = validator.ForStrings(`^[a-zA-Z]+[a-zA-Z0-9\_\.\-]*$`, 3, 50)
	uuid_        = validator.ForStrings(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, len("00000000-0000-0000-0000-000000000000"), len("00000000-0000-0000-0000-000000000000"))
)

func (v AccessId) Validate() error          { return uuid_.Validate(string(v)) }
func (v BookmarkId) Validate() error        { return uuid_.Validate(string(v)) }
func (v BottomUpPropsId) Validate() error   { return uuid_.Validate(string(v)) }
func (v CollaborationId) Validate() error   { return uuid_.Validate(string(v)) }
func (v CollaboratorId) Validate() error    { return uuid_.Validate(string(v)) }
func (v ControlAreaId) Validate() error     { return uuid_.Validate(string(v)) }
func (v DelegationId) Validate() error      { return uuid_.Validate(string(v)) }
func (v Email) Validate() error             { return email.Validate(string(v)) }
func (v GroupId) Validate() error           { return uuid_.Validate(string(v)) }
func (v GroupInviteId) Validate() error     { return uuid_.Validate(string(v)) }
func (v GroupMembershipId) Validate() error { return uuid_.Validate(string(v)) }
func (v GroupName) Validate() error         { return groupName.Validate(string(v)) }
func (v HumanName) Validate() error         { return humanName.Validate(string(v)) }
func (v LinkId) Validate() error            { return uuid_.Validate(string(v)) }
func (v LoginId) Validate() error           { return uuid_.Validate(string(v)) }
func (v ObjectiveId) Validate() error       { return uuid_.Validate(string(v)) }
func (v OperationId) Validate() error       { return uuid_.Validate(string(v)) }
func (v Phone) Validate() error             { return phoneNumber.Validate(string(v)) }
func (v PropertiesId) Validate() error      { return uuid_.Validate(string(v)) }
func (v SessionId) Validate() error         { return uuid_.Validate(string(v)) }
func (v SessionToken) Validate() error      { return sessionToken.Validate(string(v)) }
func (v TagId) Validate() error             { return uuid_.Validate(string(v)) }
func (v UserId) Validate() error            { return uuid_.Validate(string(v)) }
func (v Username) Validate() error          { return username.Validate(string(v)) }
func (v VersionId) Validate() error         { return uuid_.Validate(string(v)) }
