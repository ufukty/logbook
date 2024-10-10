package public

import (
	"fmt"

	"logbook/cmd/objectives/api/public/endpoints"
	"logbook/cmd/objectives/app"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/headers"
	"logbook/internal/web/router/cors"
	"logbook/internal/web/sidecar"
	"net/http"
	"net/url"
	"path/filepath"
)

type Public struct {
	apicfg  *api.Config
	deplcfg *deployment.Config
	e       *endpoints.Endpoints
	l       *logger.Logger
}

func New(apicfg *api.Config, deplcfg *deployment.Config, a *app.App, sc *sidecar.Sidecar, l *logger.Logger) *Public {
	l = l.Sub("Public")
	e := endpoints.New(a, l)

	return &Public{
		apicfg:  apicfg,
		deplcfg: deplcfg,
		e:       e,
		l:       l,
	}
}

func (p *Public) Register(r *http.ServeMux) error {
	s := p.apicfg.Public.Services.Objectives

	eps := map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.Attach:    p.e.ReattachObjective,
		s.Endpoints.Create:    p.e.CreateObjective,
		s.Endpoints.Mark:      p.e.MarkComplete,
		s.Endpoints.Placement: p.e.GetPlacementArray,
	}

	origin, err := url.JoinPath(p.deplcfg.Router.Cors.AllowOrigin)
	if err != nil {
		return fmt.Errorf("url.JoinPath: %w", err)
	}

	for ep, handler := range eps {
		corsed := cors.Simple(handler, origin, []string{ep.GetMethod()}, []string{headers.ContentType, headers.Authorization})
		path := filepath.Join("/public", ep.GetPath())
		for _, method := range []string{ep.GetMethod(), "OPTIONS"} {
			pattern := fmt.Sprintf("%s %s", method, path)
			p.l.Printf("registering: %s -> %p\n", pattern, corsed)
			r.HandleFunc(pattern, corsed)
		}
	}

	return nil
}
