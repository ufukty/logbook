package transports

func (a *PolicyAction) FromRoute(src string) error {
	*a = PolicyAction(src)
	return nil
}

func (a *PolicyAction) ToRoute() (string, error) {
	return string(*a), nil
}
