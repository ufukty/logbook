package columns

import (
	"slices"
)

func (v AccessId) Validate() any          { return uuid.Validate(string(v)) }
func (v BookmarkId) Validate() any        { return uuid.Validate(string(v)) }
func (v BottomUpPropsId) Validate() any   { return uuid.Validate(string(v)) }
func (v CollaborationId) Validate() any   { return uuid.Validate(string(v)) }
func (v CollaboratorId) Validate() any    { return uuid.Validate(string(v)) }
func (v ControlAreaId) Validate() any     { return uuid.Validate(string(v)) }
func (v DelegationId) Validate() any      { return uuid.Validate(string(v)) }
func (v GroupId) Validate() any           { return uuid.Validate(string(v)) }
func (v GroupInviteId) Validate() any     { return uuid.Validate(string(v)) }
func (v GroupMembershipId) Validate() any { return uuid.Validate(string(v)) }
func (v LinkId) Validate() any            { return uuid.Validate(string(v)) }
func (v LoginId) Validate() any           { return uuid.Validate(string(v)) }
func (v ObjectiveId) Validate() any       { return uuid.Validate(string(v)) }
func (v OperationId) Validate() any       { return uuid.Validate(string(v)) }
func (v PropertiesId) Validate() any      { return uuid.Validate(string(v)) }
func (v SessionId) Validate() any         { return uuid.Validate(string(v)) }
func (v TagId) Validate() any             { return uuid.Validate(string(v)) }
func (v UserId) Validate() any            { return uuid.Validate(string(v)) }
func (v VersionId) Validate() any         { return uuid.Validate(string(v)) }

func (v Email) Validate() any        { return email.Validate(string(v)) }
func (v GroupName) Validate() any    { return groupName.Validate(string(v)) }
func (v HumanName) Validate() any    { return humanName.Validate(string(v)) }
func (v Phone) Validate() any        { return phoneNumber.Validate(string(v)) }
func (v SessionToken) Validate() any { return sessionToken.Validate(string(v)) }
func (v Username) Validate() any     { return username.Validate(string(v)) }

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
