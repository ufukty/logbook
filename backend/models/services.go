package models

type Service string

const (
	Account    = Service("account")
	Auth       = Service("auth")
	Discovery  = Service("discovery")
	Groups     = Service("groups")
	Internal   = Service("internal")
	Objectives = Service("objectives")
	Tags       = Service("tags")
)

func (s *Service) Set(v string) error {
	*s = Service(v)
	return nil
}
