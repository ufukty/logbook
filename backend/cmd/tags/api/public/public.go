package public

import (
	"logbook/cmd/tags/api/public/app"
	"logbook/cmd/tags/api/public/endpoints"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/registryfile"
	"net/http"

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

func (p *Public) Endpoints() map[api.Endpoint]http.HandlerFunc {
	return map[api.Endpoint]http.HandlerFunc{
		p.apicfg.Tags.Public.Assign:   p.e.TagAssign,
		p.apicfg.Tags.Public.Creation: p.e.TagCreation,
	}
}
