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

func (ag *Agent) RegisterEndpoints(public, private map[api.Endpoint]http.HandlerFunc) error {
	origin, err := url.JoinPath(ag.deplcfg.Router.Cors.AllowOrigin)
	if err != nil {
		return fmt.Errorf("url.JoinPath: %w", err)
	}

	params := receptionistParams{
		Timeout: 1 * time.Second, // FIXME:
	}

	corsheaders := []string{
		headers.ContentType,
		headers.Authorization,
	}

	for ep, handler := range public {
		c := newCors(handler, origin, []string{ep.GetMethod()}, corsheaders)
		pl := newReceptionist(params, ag.l.Sub(ep.GetPath()), c)

		ag.l.Printf("registering: %s, OPTIONS %s -> %p\n", ep.GetMethod(), ep.GetPath(), pl)
		for _, method := range []string{ep.GetMethod(), "OPTIONS"} {
			pattern := fmt.Sprintf("%s %s", method, ep.GetPath())
			ag.r.Handle(pattern, pl)
		}
	}

	for ep, handler := range private {
		pl := newReceptionist(params, ag.l.Sub(ep.GetPath()), handler)
		pattern := fmt.Sprintf("%s %s", ep.GetMethod(), ep.GetPath())
		ag.l.Printf("registering: %s -> %p\n", pattern, pl)
		ag.r.Handle(pattern, pl)
	}

	ag.r.Handle("GET /ping", newReceptionist(params, ag.l.Sub("ping"), http.HandlerFunc(pong)))
	ag.r.Handle("GET /", newReceptionist(params, ag.l.Sub("not-found"), http.HandlerFunc(http.NotFound)))

	return nil
}

func (ag *Agent) RegisterForwarders(fwds map[string]*forwarder.LoadBalancedReverseProxy) error {
	params := receptionistParams{
		Timeout: time.Second, // FIXME:
	}

	for addr, fwd := range fwds {
		ag.l.Printf("registering forwarder for: %s -> %p\n", addr, fwd)
		l := ag.l.Sub(fmt.Sprintf("strip-prefix(%s)", addr))
		ag.r.Handle(addr+"/", newReceptionist(params, l, http.StripPrefix(addr, fwd)))
	}

	ag.r.Handle("/ping", newReceptionist(params, ag.l.Sub("ping"), http.HandlerFunc(pong)))
	ag.r.Handle("/", newReceptionist(params, ag.l.Sub("not-found"), http.HandlerFunc(http.NotFound)))

	return nil
}
