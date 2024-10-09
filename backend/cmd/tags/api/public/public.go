package public

import (
	"fmt"
	"logbook/cmd/tags/api/public/app"
	"logbook/cmd/tags/api/public/endpoints"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/headers"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router/cors"
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

func New(apicfg *api.Config, deplcfg *deployment.Config, pool *pgxpool.Pool, internalsd *registryfile.FileReader, l *logger.Logger) *Public {
	a := app.New(pool, apicfg, internalsd)
	e := endpoints.New(a, l)

	return &Public{
		pool:   pool,
		apicfg: apicfg,
		e:      e,
	}
}

func (p *Public) Register(r *http.ServeMux) error {
	s := p.apicfg.Public.Services.Tags

	eps := map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.Assign:   p.e.TagAssign,
		s.Endpoints.Creation: p.e.TagCreation,
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
