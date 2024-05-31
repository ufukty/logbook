package models

type Service string

var (
	Account    = Service("account")
	Gateway    = Service("gateway")
	Objectives = Service("objectives")
)
