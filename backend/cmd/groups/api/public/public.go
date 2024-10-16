package public

import (
	"logbook/cmd/groups/api/public/app"
	"logbook/cmd/groups/api/public/endpoints"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/sidecar"
	"net/http"

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

func (p *Public) Endpoints() map[api.Endpoint]http.HandlerFunc {
	return map[api.Endpoint]http.HandlerFunc{
		p.apicfg.Groups.Public.Create: p.em.CreateGroup,
	}
}
