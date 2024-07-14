package models

import "fmt"

type Service string

var (
	Account    = Service("account")
	Discovery  = Service("discovery")
	Internal   = Service("internal")
	Objectives = Service("objectives")
)

func (s *Service) Set(v string) error {
	switch v {
	case string(Account),
		string(Discovery),
		string(Internal),
		string(Objectives):
		*s = Service(v)
	default:
		return fmt.Errorf("invalid service name: %s", v)
	}
	return nil
}
