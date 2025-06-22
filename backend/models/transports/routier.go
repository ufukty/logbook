package transports

func (a *Country) FromRoute(src string) error        { *a = Country(src); return nil }
func (a *InviteResponse) FromRoute(src string) error { *a = InviteResponse(src); return nil }
func (a *MemberType) FromRoute(src string) error     { *a = MemberType(src); return nil }
func (a *PolicyAction) FromRoute(src string) error   { *a = PolicyAction(src); return nil }

func (a *Country) ToRoute() (string, error)        { return string(*a), nil }
func (a *InviteResponse) ToRoute() (string, error) { return string(*a), nil }
func (a *MemberType) ToRoute() (string, error)     { return string(*a), nil }
func (a *PolicyAction) ToRoute() (string, error)   { return string(*a), nil }
