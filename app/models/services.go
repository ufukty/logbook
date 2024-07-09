package models

type Service string

var (
	Account    = Service("account")
	Discovery  = Service("discovery")
	Internal   = Service("internal")
	Objectives = Service("objectives")
)
