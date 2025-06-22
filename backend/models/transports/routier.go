package transports

func (a *PolicyAction) FromRoute(src string) error {
	*a = PolicyAction(src)
	err := a.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (a *PolicyAction) ToRoute() (string, error) {
	return string(*a), nil
}
