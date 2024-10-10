package public

import (
	"fmt"
	"logbook/cmd/account/api/public/app"
	"logbook/cmd/account/api/public/endpoints"
	objectives "logbook/cmd/objectives/api/private/client"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/balancer"
	"logbook/internal/web/headers"
	"logbook/internal/web/router/cors"
	"logbook/internal/web/sidecar"
	"logbook/models"
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

func (p *Public) Register(r *http.ServeMux) error {
	s := p.apicfg.Public.Services.Account

	eps := map[api.Endpoint]http.HandlerFunc{
		s.Endpoints.CreateAccount: p.em.CreateAccount,
		s.Endpoints.CreateProfile: p.em.CreateProfile,
		s.Endpoints.Login:         p.em.Login,
		s.Endpoints.Logout:        p.em.Logout,
		s.Endpoints.Whoami:        p.em.WhoAmI,
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
