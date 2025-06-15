package columns

import (
	"github.com/ufukty/gohandlers/pkg/validator/validate"
)

func (v NonNegativeNumber) Validate() error {
	if v < 0 {
		return validate.ErrPattern
	}
	return nil
}

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
