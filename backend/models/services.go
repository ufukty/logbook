package models

type Service string

const (
	Auth         Service = "auth"
	Discovery    Service = "discovery"
	Groups       Service = "groups"
	Internal     Service = "internal"
	Objectives   Service = "objectives"
	Pdp          Service = "pdp"
	Profiles     Service = "profiles"
	Registration Service = "registration"
	Registry     Service = "registry"
	Sessions     Service = "sessions"
	Tags         Service = "tags"
	Users        Service = "users"
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
