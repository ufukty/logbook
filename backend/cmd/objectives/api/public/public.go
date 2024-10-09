package public

import (
	"fmt"
	"logbook/cmd/objectives/api/public/app"
	"logbook/cmd/objectives/api/public/endpoints"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/headers"
	"logbook/internal/web/router/cors"
	"logbook/internal/web/sidecar"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Public struct {
	pool    *pgxpool.Pool
	apicfg  *api.Config
	deplcfg *deployment.Config
	e       *endpoints.Endpoints
}

func New(apicfg *api.Config, deplcfg *deployment.Config, pool *pgxpool.Pool, sc *sidecar.Sidecar, l *logger.Logger) *Public {
	a := app.New(pool, l)
	e := endpoints.New(a, l)

	return &Public{
		pool:   pool,
		apicfg: apicfg,
		e:      e,
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
		r.HandleFunc(fmt.Sprintf("OPTIONS %s", path), corsed)
		r.HandleFunc(fmt.Sprintf("%s %s", ep.GetMethod(), path), corsed)
	}

	return nil
}
