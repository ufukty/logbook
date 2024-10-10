package public

import (
	"fmt"
	"logbook/cmd/groups/api/public/app"
	"logbook/cmd/groups/api/public/endpoints"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/headers"
	"logbook/internal/web/router/cors"
	"logbook/internal/web/sidecar"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Public struct {
	pool    *pgxpool.Pool
	apicfg  *api.Config
	deplcfg *deployment.Config
	em      *endpoints.Endpoints
	l       *logger.Logger
}

func New(apicfg *api.Config, deplcfg *deployment.Config, pool *pgxpool.Pool, sc *sidecar.Sidecar, l *logger.Logger) *Public {
	l = l.Sub("Public")

	a := app.New(pool)
	e := endpoints.New(a, l)

	return &Public{
		pool:    pool,
		apicfg:  apicfg,
		deplcfg: deplcfg,
		em:      e,
		l:       l,
	}
}

func (p *Public) Register(r *http.ServeMux) error {
	s := p.apicfg.Public.Services.Groups

	eps := map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.Create: p.em.CreateGroup,
	}

	origin, err := url.JoinPath(p.deplcfg.Router.Cors.AllowOrigin)
	if err != nil {
		return fmt.Errorf("url.JoinPath: %w", err)
	}

	for ep, handler := range eps {
		corsed := cors.Simple(handler, origin, []string{ep.GetMethod()}, []string{headers.ContentType, headers.Authorization})
		path := api.ByService(ep)
		for _, method := range []string{ep.GetMethod(), "OPTIONS"} {
			pattern := fmt.Sprintf("%s %s", method, path)
			p.l.Printf("registering: %s -> %p\n", pattern, corsed)
			r.HandleFunc(pattern, corsed)
		}
	}

	return nil
}
