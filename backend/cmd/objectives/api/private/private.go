package private

import (
	"fmt"
	"logbook/cmd/objectives/api/private/endpoints"
	"logbook/cmd/objectives/app"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"net/http"
)

type Private struct {
	apicfg  *api.Config
	em      *endpoints.Endpoints
	deplcfg *deployment.Config
	l       *logger.Logger
}

func New(apicfg *api.Config, deplcfg *deployment.Config, a *app.App, l *logger.Logger) *Private {
	l = l.Sub("Private")

	em := endpoints.New(a, l)

	return &Private{
		apicfg:  apicfg,
		em:      em,
		deplcfg: deplcfg,
		l:       l,
	}
}

func (p *Private) Register(r *http.ServeMux) error {
	s := p.apicfg.Internal.Services.Objectives

	eps := map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.RockCreate: p.em.RockCreate,
	}

	for ep, handler := range eps {
		pattern := fmt.Sprintf("%s %s", ep.GetMethod(), api.ByService(ep))
		p.l.Printf("registering: %s -> %p\n", pattern, handler)
		r.HandleFunc(pattern, handler)
	}

	return nil
}
