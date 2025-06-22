package columns

func (v *AccessId) FromRoute(src string) error          { *v = AccessId(src); return nil }
func (v *BookmarkId) FromRoute(src string) error        { *v = BookmarkId(src); return nil }
func (v *BottomUpPropsId) FromRoute(src string) error   { *v = BottomUpPropsId(src); return nil }
func (v *CollaborationId) FromRoute(src string) error   { *v = CollaborationId(src); return nil }
func (v *CollaboratorId) FromRoute(src string) error    { *v = CollaboratorId(src); return nil }
func (v *ControlAreaId) FromRoute(src string) error     { *v = ControlAreaId(src); return nil }
func (v *DelegationId) FromRoute(src string) error      { *v = DelegationId(src); return nil }
func (v *Email) FromRoute(src string) error             { *v = Email(src); return nil }
func (v *GroupId) FromRoute(src string) error           { *v = GroupId(src); return nil }
func (v *GroupInviteId) FromRoute(src string) error     { *v = GroupInviteId(src); return nil }
func (v *GroupMembershipId) FromRoute(src string) error { *v = GroupMembershipId(src); return nil }
func (v *GroupName) FromRoute(src string) error         { *v = GroupName(src); return nil }
func (v *HumanName) FromRoute(src string) error         { *v = HumanName(src); return nil }
func (v *LinkId) FromRoute(src string) error            { *v = LinkId(src); return nil }
func (v *LinkType) FromRoute(src string) error          { *v = LinkType(src); return nil }
func (v *LoginId) FromRoute(src string) error           { *v = LoginId(src); return nil }
func (v *ObjectiveId) FromRoute(src string) error       { *v = ObjectiveId(src); return nil }
func (v *OperationId) FromRoute(src string) error       { *v = OperationId(src); return nil }
func (v *PropertiesId) FromRoute(src string) error      { *v = PropertiesId(src); return nil }
func (v *SessionId) FromRoute(src string) error         { *v = SessionId(src); return nil }
func (v *SessionToken) FromRoute(src string) error      { *v = SessionToken(src); return nil }
func (v *TagId) FromRoute(src string) error             { *v = TagId(src); return nil }
func (v *TagTitle) FromRoute(src string) error          { *v = TagTitle(src); return nil }
func (v *UserId) FromRoute(src string) error            { *v = UserId(src); return nil }
func (v *Username) FromRoute(src string) error          { *v = Username(src); return nil }
func (v *VersionId) FromRoute(src string) error         { *v = VersionId(src); return nil }

func (v AccessId) ToRoute() (string, error)          { return string(v), nil }
func (v BookmarkId) ToRoute() (string, error)        { return string(v), nil }
func (v BottomUpPropsId) ToRoute() (string, error)   { return string(v), nil }
func (v CollaborationId) ToRoute() (string, error)   { return string(v), nil }
func (v CollaboratorId) ToRoute() (string, error)    { return string(v), nil }
func (v ControlAreaId) ToRoute() (string, error)     { return string(v), nil }
func (v DelegationId) ToRoute() (string, error)      { return string(v), nil }
func (v Email) ToRoute() (string, error)             { return string(v), nil }
func (v GroupId) ToRoute() (string, error)           { return string(v), nil }
func (v GroupInviteId) ToRoute() (string, error)     { return string(v), nil }
func (v GroupMembershipId) ToRoute() (string, error) { return string(v), nil }
func (v GroupName) ToRoute() (string, error)         { return string(v), nil }
func (v HumanName) ToRoute() (string, error)         { return string(v), nil }
func (v LinkId) ToRoute() (string, error)            { return string(v), nil }
func (v LinkType) ToRoute() (string, error)          { return string(v), nil }
func (v LoginId) ToRoute() (string, error)           { return string(v), nil }
func (v ObjectiveId) ToRoute() (string, error)       { return string(v), nil }
func (v OperationId) ToRoute() (string, error)       { return string(v), nil }
func (v PropertiesId) ToRoute() (string, error)      { return string(v), nil }
func (v SessionId) ToRoute() (string, error)         { return string(v), nil }
func (v SessionToken) ToRoute() (string, error)      { return string(v), nil }
func (v TagId) ToRoute() (string, error)             { return string(v), nil }
func (v TagTitle) ToRoute() (string, error)          { return string(v), nil }
func (v UserId) ToRoute() (string, error)            { return string(v), nil }
func (v Username) ToRoute() (string, error)          { return string(v), nil }
func (v VersionId) ToRoute() (string, error)         { return string(v), nil }
