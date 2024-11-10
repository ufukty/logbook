package models

type Service string

const (
	Account    Service = "account"
	Auth       Service = "auth"
	Discovery  Service = "discovery"
	Groups     Service = "groups"
	Internal   Service = "internal"
	Objectives Service = "objectives"
	Registry   Service = "registry"
	Tags       Service = "tags"
)

func (s *Service) Set(v string) error {
	*s = Service(v)
	return nil
}

func (s *Service) FromRoute(src string) error {
	*s = Service(src)
	return nil
}

func (s Service) ToRoute() (string, error) {
	return string(s), nil
}
