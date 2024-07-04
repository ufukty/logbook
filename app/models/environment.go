package models

type Environment string

const (
	Local      = Environment("local")
	Stage      = Environment("stage")
	Production = Environment("production")
)
