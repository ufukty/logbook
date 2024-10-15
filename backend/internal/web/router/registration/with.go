package registration

import (
	"fmt"
	"logbook/config/api"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/headers"
	"logbook/internal/web/router/registration/middlewares"
	"logbook/internal/web/router/registration/receptionist"
	"logbook/internal/web/router/registration/receptionist/decls"
	"net/http"
	"net/url"
	"time"
)

// TODO: auth
func (ag *Agent) RegisterForInternal(eps map[api.Endpoint]decls.HandlerFunc) error {
	var (
		a = middlewares.NewAuth()
	)

	params := receptionist.Params{
		Timeout: time.Second, // FIXME:
	}

	for ep, handler := range eps {
		pl := receptionist.New(params, ag.l.Sub(api.ByService(ep)),
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
func (ag *Agent) RegisterForPublic(eps map[api.Endpoint]decls.HandlerFunc) error {
	origin, err := url.JoinPath(ag.deplcfg.Router.Cors.AllowOrigin)
	if err != nil {
		return fmt.Errorf("url.JoinPath: %w", err)
	}

	var (
		a  = middlewares.NewAuth()
		cm = middlewares.NewCorsManager(origin)
	)

	params := receptionist.Params{
		Timeout: 1 * time.Second, // FIXME:
	}
	corsheaders := []string{
		headers.ContentType,
		headers.Authorization,
	}

	for ep, handler := range eps {
		pl := receptionist.New(params, ag.l.Sub(api.ByService(ep)),
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
	params := receptionist.Params{
		Timeout: time.Second, // FIXME:
	}

	for addr, fwd := range fwds {
		str := middlewares.NewStripper(servicepath, fwd)
		route := api.PrefixedByGateway(addr)
		ag.r.Handle(route, receptionist.New(params, ag.l.Sub(fmt.Sprintf("Stipper(%s)", route)), str.Strip))
	}

	return nil
}

func (ag *Agent) RegisterCommonalities() error {
	params := receptionist.Params{
		Timeout: 1 * time.Second, // FIXME:
	}

	ag.r.Handle("/ping", receptionist.New(params, ag.l.Sub("ping"), middlewares.Pong))
	ag.r.Handle("/", http.HandlerFunc(http.NotFound))

	return nil
}
