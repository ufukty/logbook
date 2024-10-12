package public

import (
	"fmt"
	"time"

	"logbook/cmd/objectives/api/public/endpoints"
	"logbook/cmd/objectives/app"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/headers"
	"logbook/internal/web/router/pipelines"
	"logbook/internal/web/router/pipelines/middlewares"

	"logbook/internal/web/sidecar"
	"net/http"
	"net/url"
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

	eps := map[api.Endpoint]pipelines.HandlerFunc[middlewares.Store]{
		s.Endpoints.Attach:    p.e.ReattachObjective,
		s.Endpoints.Create:    p.e.CreateObjective,
		s.Endpoints.Mark:      p.e.MarkComplete,
		s.Endpoints.Placement: p.e.GetPlacementArray,
	}

	origin, err := url.JoinPath(p.deplcfg.Router.Cors.AllowOrigin)
	if err != nil {
		return fmt.Errorf("url.JoinPath: %w", err)
	}

	// TODO: log
	// TODO: not found
	// DONE: recover
	// DONE: cors
	// TODO: auth
	// TODO: timeout
	// TODO: *handler
	var (
		a  = middlewares.NewAuth()
		cm = middlewares.NewCorsManager(origin)
	)

	for ep, handler := range eps {
		pl := pipelines.NewPipeline([]pipelines.HandlerFunc[middlewares.Store]{
			a.Handle,
			cm.Instantiate([]string{ep.GetMethod()}, []string{headers.ContentType, headers.Authorization}).Handle,
			handler,
		}, pipelines.PipelineParams{
			Timeout: 1 * time.Second,
		}, p.l.Sub(api.ByService(ep)))
		for _, method := range []string{ep.GetMethod(), "OPTIONS"} {
			pattern := fmt.Sprintf("%s %s", method, api.ByService(ep))
			p.l.Printf("registering: %s -> %p\n", pattern, pl)
			r.Handle(pattern, pl)
		}
	}

	return nil
}
