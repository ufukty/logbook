package public

import (
	"logbook/cmd/account/api/public/app"
	"logbook/cmd/account/api/public/endpoints"
	objectives "logbook/cmd/objectives/api/private/client"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/balancer"
	"logbook/internal/web/router/receptionist"
	"logbook/internal/web/router/registration"
	"logbook/internal/web/router/registration/middlewares"
	"logbook/internal/web/sidecar"
	"logbook/models"

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

	objectives := objectives.NewClient(balancer.New(sc.InstanceSource(models.Objectives)), apicfg)
	app := app.New(pool, apicfg, objectives)
	em := endpoints.New(app, l)

	return &Public{
		pool:    pool,
		apicfg:  apicfg,
		deplcfg: deplcfg,
		em:      em,
		l:       l,
	}
}

func (p *Public) Register(agent *registration.Agent) error {
	s := p.apicfg.Public.Services.Account
	return agent.RegisterForPublic(map[api.Endpoint]receptionist.HandlerFunc[middlewares.Store]{
		s.Endpoints.CreateAccount: p.em.CreateAccount,
		s.Endpoints.CreateProfile: p.em.CreateProfile,
		s.Endpoints.Login:         p.em.Login,
		s.Endpoints.Logout:        p.em.Logout,
		s.Endpoints.Whoami:        p.em.WhoAmI,
	})
}
