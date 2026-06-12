package register

import (
	"fmt"
	"net/http"
	"net/url"

	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/registration/reception"
	"logbook/models"

	"go.ufukty.com/gohandlers/pkg/gohandlers"
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

type Lister interface {
	ListHandlers() map[string]gohandlers.HandlerInfo
}

func (ag *Agent) RegisterEndpoints(public, private Lister) error {
	origin, err := url.JoinPath(ag.deplcfg.Router.Cors.AllowOrigin)
	if err != nil {
		return fmt.Errorf("url.JoinPath: %w", err)
	}

	corsheaders := []string{"Authorization", "Content-Type"}

	if public != nil {
		for hn, info := range public.ListHandlers() {
			c := reception.NewCors(info.Ref, origin, []string{info.Method}, corsheaders)
			pl := reception.New(ag.deplcfg, ag.l.Sub(info.Path), c)

			ag.l.Printf("registering: %s (%s, OPTIONS %s) -> %p\n", hn, info.Method, info.Path, pl)
			for _, method := range []string{info.Method, "OPTIONS"} {
				pattern := fmt.Sprintf("%s %s", method, info.Path)
				ag.r.Handle(pattern, pl)
			}
		}
	}

	if private != nil {
		for hn, info := range private.ListHandlers() {
			pl := reception.New(ag.deplcfg, ag.l.Sub(info.Path), info.Ref)
			pattern := fmt.Sprintf("%s %s", info.Method, info.Path)
			ag.l.Printf("registering: %s (%s) -> %p\n", hn, pattern, pl)
			ag.r.Handle(pattern, pl)
		}
	}

	ag.r.Handle("GET /ping", reception.New(ag.deplcfg, ag.l.Sub("ping"), http.HandlerFunc(reception.Pong)))
	ag.r.Handle("GET /", reception.New(ag.deplcfg, ag.l.Sub("not-found"), http.HandlerFunc(http.NotFound)))

	return nil
}

func (ag *Agent) RegisterForwarders(fwds map[models.Service]*forwarder.LoadBalancedReverseProxy) error {
	for addr, fwd := range fwds {
		ag.l.Printf("registering forwarder for: %s -> %p\n", addr, fwd)
		l := ag.l.Sub(fmt.Sprintf("strip-prefix(%s)", addr))
		ag.r.Handle(string(addr)+"/", reception.New(ag.deplcfg, l, http.StripPrefix(string(addr), fwd)))
	}

	ag.r.Handle("/ping", reception.New(ag.deplcfg, ag.l.Sub("ping"), http.HandlerFunc(reception.Pong)))
	ag.r.Handle("/", reception.New(ag.deplcfg, ag.l.Sub("not-found"), http.HandlerFunc(http.NotFound)))

	return nil
}
