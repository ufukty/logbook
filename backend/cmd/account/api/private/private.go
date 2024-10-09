package private

import (
	"fmt"
	"logbook/cmd/account/api/private/app"
	"logbook/cmd/account/api/private/endpoints"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"net/http"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Private struct {
	pool    *pgxpool.Pool
	apicfg  *api.Config
	deplcfg *deployment.Config
	em      *endpoints.Endpoints
}

func New(apicfg *api.Config, deplcfg *deployment.Config, pool *pgxpool.Pool, l *logger.Logger) *Private {
	app := app.New(pool, apicfg)
	em := endpoints.New(app, l)

	return &Private{
		pool:   pool,
		apicfg: apicfg,
		em:     em,
	}
}

func (p *Private) Register(r *http.ServeMux) error {
	s := p.apicfg.Internal.Services.Account

	eps := map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.WhoIs: p.em.WhoIs,
	}

	for ep, handler := range eps {
		path := filepath.Join("/private", ep.GetPath())
		r.HandleFunc(fmt.Sprintf("%s %s", ep.GetMethod(), path), handler)
	}

	return nil
}
