package private

import (
	"logbook/cmd/account/api/private/app"
	"logbook/cmd/account/api/private/endpoints"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/router/registration"
	"logbook/internal/web/router/registration/receptionist/decls"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Private struct {
	pool    *pgxpool.Pool
	apicfg  *api.Config
	deplcfg *deployment.Config
	em      *endpoints.Endpoints
	l       *logger.Logger
}

func New(apicfg *api.Config, deplcfg *deployment.Config, pool *pgxpool.Pool, l *logger.Logger) *Private {
	l = l.Sub("Private")
	app := app.New(pool, apicfg)
	em := endpoints.New(app, l)

	return &Private{
		pool:    pool,
		apicfg:  apicfg,
		deplcfg: deplcfg,
		em:      em,
		l:       l,
	}
}

func (p *Private) Register(agent *registration.Agent) error {
	s := p.apicfg.Internal.Services.Account
	return agent.RegisterForInternal(map[api.Endpoint]decls.HandlerFunc{
		s.Endpoints.WhoIs: p.em.WhoIs,
	})
}
