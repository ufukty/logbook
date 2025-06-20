package columns

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

