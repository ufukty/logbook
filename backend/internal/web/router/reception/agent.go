package reception

import (
	"fmt"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/headers"
	"net/http"
	"net/url"
	"time"
)

// [Agent] is the registration Agent which helps services, and gateways to register their handlers and forwarders appropriately
type Agent struct {
	deplcfg *deployment.Config
	r       *http.ServeMux
	l       *logger.Logger
}

func NewAgent(deplcfg *deployment.Config, l *logger.Logger) *Agent {
	return &Agent{
		deplcfg: deplcfg,
		r:       http.NewServeMux(),
		l:       l.Sub("agent"),
	}
}

func (a *Agent) Mux() *http.ServeMux {
	return a.r
}

// TODO: auth
func (ag *Agent) RegisterForInternal(eps map[api.Endpoint]HandlerFunc) error {
	var (
		a = NewAuth()
	)

	params := receptionistParams{
		Timeout: time.Second, // FIXME:
	}

	for ep, handler := range eps {
		pl := newReceptionist(params, ag.l.Sub(api.ByService(ep)),
			a.Handle,
			handler,
		)
		pattern := fmt.Sprintf("%s %s", ep.GetMethod(), api.ByService(ep))
		ag.l.Printf("registering: %s -> %p\n", pattern, handler)
		ag.r.Handle(pattern, pl)
	}

	return nil
}

// DONE: cors
// TODO: auth
func (ag *Agent) RegisterForPublic(eps map[api.Endpoint]HandlerFunc) error {
	origin, err := url.JoinPath(ag.deplcfg.Router.Cors.AllowOrigin)
	if err != nil {
		return fmt.Errorf("url.JoinPath: %w", err)
	}

	var (
		a  = NewAuth()
		cm = NewCorsManager(origin)
	)

	params := receptionistParams{
		Timeout: 1 * time.Second, // FIXME:
	}
	corsheaders := []string{
		headers.ContentType,
		headers.Authorization,
	}

	for ep, handler := range eps {
		pl := newReceptionist(params, ag.l.Sub(api.ByService(ep)),
			a.Handle,
			cm.Instantiate([]string{ep.GetMethod()}, corsheaders).Handle,
			handler,
		)
		for _, method := range []string{ep.GetMethod(), "OPTIONS"} {
			pattern := fmt.Sprintf("%s %s", method, api.ByService(ep))
			ag.l.Printf("registering: %s -> %p\n", pattern, handler)
			ag.r.Handle(pattern, pl)
		}
	}

	return nil
}

func (ag *Agent) RegisterForwarders(servicepath string, fwds map[api.Addressable]*forwarder.LoadBalancedReverseProxy) error {
	params := receptionistParams{
		Timeout: time.Second, // FIXME:
	}

	for addr, fwd := range fwds {
		str := NewStripper(servicepath, fwd)
		route := api.PrefixedByGateway(addr)
		ag.r.Handle(route, newReceptionist(params, ag.l.Sub(fmt.Sprintf("stripper(%s)", route)), str.Strip))
	}

	return nil
}

func (ag *Agent) RegisterCommonalities() error {
	params := receptionistParams{
		Timeout: 1 * time.Second, // FIXME:
	}

	ag.r.Handle("/ping", newReceptionist(params, ag.l.Sub("ping"), Pong))
	ag.r.Handle("/", newReceptionist(params, ag.l.Sub("not found"), NotFound))

	return nil
}
