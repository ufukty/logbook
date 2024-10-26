package public

import (
	"logbook/cmd/account/api/public/app"
	"logbook/cmd/account/api/public/endpoints"
	objectives "logbook/cmd/objectives/client"
	"logbook/config/api"
	"logbook/internal/logger"
	"logbook/internal/web/balancer"
	"logbook/internal/web/sidecar"
	"logbook/models"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Public struct {
	pool   *pgxpool.Pool
	apicfg *api.Config
	em     *endpoints.Endpoints
	l      *logger.Logger
}

func New(apicfg *api.Config, pool *pgxpool.Pool, sc *sidecar.Sidecar, l *logger.Logger) *Public {
	l = l.Sub("Public")

	objectives := objectives.NewClient(balancer.New(sc.InstanceSource(models.Objectives)), apicfg)
	app := app.New(pool, apicfg, objectives)
	em := endpoints.New(app, l)

	return &Public{
		pool:   pool,
		apicfg: apicfg,
		em:     em,
		l:      l,
	}
}

func (p *Public) Endpoints() map[api.Endpoint]http.HandlerFunc {
	s := p.apicfg.Account.Public
	return map[api.Endpoint]http.HandlerFunc{
		s.CreateAccount: p.em.CreateAccount,
		s.CreateProfile: p.em.CreateProfile,
		s.Login:         p.em.Login,
		s.Logout:        p.em.Logout,
		s.Whoami:        p.em.WhoAmI,
	}
}
