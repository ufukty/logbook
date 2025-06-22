package columns

import (
	"logbook/models/validators"
	"slices"
)

func (v AccessId) Validate() any          { return validators.Uuid.Validate(string(v)) }
func (v BookmarkId) Validate() any        { return validators.Uuid.Validate(string(v)) }
func (v BottomUpPropsId) Validate() any   { return validators.Uuid.Validate(string(v)) }
func (v CollaborationId) Validate() any   { return validators.Uuid.Validate(string(v)) }
func (v CollaboratorId) Validate() any    { return validators.Uuid.Validate(string(v)) }
func (v ControlAreaId) Validate() any     { return validators.Uuid.Validate(string(v)) }
func (v DelegationId) Validate() any      { return validators.Uuid.Validate(string(v)) }
func (v GroupId) Validate() any           { return validators.Uuid.Validate(string(v)) }
func (v GroupInviteId) Validate() any     { return validators.Uuid.Validate(string(v)) }
func (v GroupMembershipId) Validate() any { return validators.Uuid.Validate(string(v)) }
func (v LinkId) Validate() any            { return validators.Uuid.Validate(string(v)) }
func (v LoginId) Validate() any           { return validators.Uuid.Validate(string(v)) }
func (v ObjectiveId) Validate() any       { return validators.Uuid.Validate(string(v)) }
func (v OperationId) Validate() any       { return validators.Uuid.Validate(string(v)) }
func (v PropertiesId) Validate() any      { return validators.Uuid.Validate(string(v)) }
func (v SessionId) Validate() any         { return validators.Uuid.Validate(string(v)) }
func (v TagId) Validate() any             { return validators.Uuid.Validate(string(v)) }
func (v UserId) Validate() any            { return validators.Uuid.Validate(string(v)) }
func (v VersionId) Validate() any         { return validators.Uuid.Validate(string(v)) }

func (v Email) Validate() any        { return validators.Email.Validate(string(v)) }
func (v GroupName) Validate() any    { return validators.GroupName.Validate(string(v)) }
func (v HumanName) Validate() any    { return validators.HumanName.Validate(string(v)) }
func (v Phone) Validate() any        { return validators.PhoneNumber.Validate(string(v)) }
func (v SessionToken) Validate() any { return validators.SessionToken.Validate(string(v)) }
func (v Username) Validate() any     { return validators.Username.Validate(string(v)) }

func (v ObjectiveContent) Validate() any { return nil }
func (v TagTitle) Validate() any         { return nil }

func collect[E interface{ Validate() any }](s []E) any {
	issues := []any{}
	for _, e := range s {
		if issue := e.Validate(); issue != nil {
			issues = append(issues, issue)
		}
	}
	if len(issues) > 0 {
		return issues
	}
	return nil
}

func (gs GroupIds) Validate() any { return collect(gs) }
func (us UserIds) Validate() any  { return collect(us) }

func (l LinkType) Validate() any {
	if !slices.Contains([]LinkType{Primary, Private, Remote}, l) {
		return "invalid value"
	}
	return nil
}
