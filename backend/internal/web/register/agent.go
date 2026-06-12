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

type Lister interface {
	ListHandlers() map[string]gohandlers.HandlerInfo
}

func debug(c *deployment.Config, l *logger.Logger, r *http.ServeMux) {
	r.Handle("GET /ping", reception.New(c, l.Sub("ping"), http.HandlerFunc(reception.Pong)))
	r.Handle("GET /", reception.New(c, l.Sub("not-found"), http.HandlerFunc(http.NotFound)))
}

// For non-gateway services
func Endpoints(c *deployment.Config, l *logger.Logger, public, private Lister) (*http.ServeMux, error) {
	r := http.NewServeMux()
	l = l.Sub("register")

	origin, err := url.JoinPath(c.Router.Cors.AllowOrigin)
	if err != nil {
		return nil, fmt.Errorf("url.JoinPath: %w", err)
	}

	corsheaders := []string{"Authorization", "Content-Type"}

	if public != nil {
		for hn, info := range public.ListHandlers() {
			cors := reception.NewCors(info.Ref, origin, []string{info.Method}, corsheaders)
			pl := reception.New(c, l.Sub(info.Path), cors)

			l.Printf("registering: %s (%s, OPTIONS %s) -> %p\n", hn, info.Method, info.Path, pl)
			for _, method := range []string{info.Method, "OPTIONS"} {
				pattern := fmt.Sprintf("%s %s", method, info.Path)
				r.Handle(pattern, pl)
			}
		}
	}

	if private != nil {
		for hn, info := range private.ListHandlers() {
			pl := reception.New(c, l.Sub(info.Path), info.Ref)
			pattern := fmt.Sprintf("%s %s", info.Method, info.Path)
			l.Printf("registering: %s (%s) -> %p\n", hn, pattern, pl)
			r.Handle(pattern, pl)
		}
	}

	debug(c, l, r)

	return r, nil
}

func Forwarders(c *deployment.Config, l *logger.Logger, fwds map[models.Service]*forwarder.LoadBalancedReverseProxy) (*http.ServeMux, error) {
	r := http.NewServeMux()
	l = l.Sub("register")

	for addr, fwd := range fwds {
		l.Printf("registering forwarder for: %s -> %p\n", addr, fwd)
		l := l.Sub(fmt.Sprintf("strip-prefix(%s)", addr))
		r.Handle(string(addr)+"/", reception.New(c, l, http.StripPrefix(string(addr), fwd)))
	}

	debug(c, l, r)

	return r, nil
}
