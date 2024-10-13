package public

import (
	"logbook/cmd/tags/api/public/app"
	"logbook/cmd/tags/api/public/endpoints"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/router/receptionist"
	"logbook/internal/web/router/registration"
	"logbook/internal/web/router/registration/middlewares"

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

func (p *Public) Register(agent *registration.Agent) error {
	s := p.apicfg.Public.Services.Tags
	return agent.RegisterForPublic(map[api.Endpoint]receptionist.HandlerFunc[middlewares.Store]{
		s.Endpoints.Assign:   p.e.TagAssign,
		s.Endpoints.Creation: p.e.TagCreation,
	})
}
