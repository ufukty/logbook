package private

import (
	"fmt"
	"logbook/cmd/objectives/api/private/endpoints"
	"logbook/cmd/objectives/app"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"net/http"
	"path/filepath"
)

type Private struct {
	apicfg *api.Config
	em     *endpoints.Endpoints
}

func New(apicfg *api.Config, deplcfg *deployment.Config, a *app.App, l *logger.Logger) *Private {
	em := endpoints.New(a, l)

	return &Private{
		apicfg: apicfg,
		em:     em,
	}
}

func (p *Private) Register(r *http.ServeMux) error {
	s := p.apicfg.Internal.Services.Objectives

	eps := map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.RockCreate: p.em.RockCreate,
	}

	for ep, handler := range eps {
		path := filepath.Join("/private", ep.GetPath())
		r.HandleFunc(fmt.Sprintf("%s %s", ep.GetMethod(), path), handler)
	}

	return nil
}
