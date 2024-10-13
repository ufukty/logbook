package registration

import (
	"logbook/config/deployment"
	"logbook/internal/logger"
	"net/http"
)

type Agent struct {
	deplcfg *deployment.Config
	r       *http.ServeMux
	l       *logger.Logger
}

func New(deplcfg *deployment.Config, l *logger.Logger) *Agent {
	return &Agent{
		deplcfg: deplcfg,
		r:       http.NewServeMux(),
		l:       l.Sub("registration"),
	}
}

func (a *Agent) Mux() *http.ServeMux {
	return a.r
}
