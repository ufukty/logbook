package models

import "fmt"

type Service string

var (
	Account    = Service("account")
	Discovery  = Service("discovery")
	Groups     = Service("groups")
	Internal   = Service("internal")
	Objectives = Service("objectives")
	Tags       = Service("tags")
)

func (s *Service) Set(v string) error {
	switch v {
	case
		string(Account),
		string(Discovery),
		string(Groups),
		string(Internal),
		string(Objectives),
		string(Tags):
		*s = Service(v)
	default:
		return fmt.Errorf("invalid service name: %s", v)
	}
	return nil
}
