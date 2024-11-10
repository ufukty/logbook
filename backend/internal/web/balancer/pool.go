package balancer

type Pool interface {
	Host() (string, error)
}
