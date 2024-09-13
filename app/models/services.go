package models

import "fmt"

type Service string

var (
	Account    = Service("account")
	Discovery  = Service("discovery")
	Internal   = Service("internal")
	Objectives = Service("objectives")
	Groups     = Service("groups")
)

func (s *Service) Set(v string) error {
	switch v {
	case string(Account),
		string(Discovery),
		string(Internal),
		string(Objectives),
		string(Groups):
		*s = Service(v)
	default:
		return fmt.Errorf("invalid service name: %s", v)
	}
	return nil
}
