package public

import (
	"logbook/cmd/groups/api/public/app"
	"logbook/cmd/groups/api/public/endpoints"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/router/reception"
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

func (p *Public) Register(agent *reception.Agent) error {
	s := p.apicfg.Public.Services.Groups
	return agent.RegisterForPublic(map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.Create: p.em.CreateGroup,
	})
	return nil
}
