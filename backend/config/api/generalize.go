package api

type Endpoint interface {
	GetPath() string
	GetMethod() string
}
