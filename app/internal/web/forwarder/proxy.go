package forwarder

import (
	"fmt"
	"log"
	"logbook/internal/web/balancer"
	"logbook/internal/web/discovery"
	"logbook/internal/web/logger"
	"logbook/models"
	"net/http"
	"net/http/httputil"
	"strings"
)

// TODO: load balancing between processes listen different ports but in same IP address
type LoadBalancedReverseProxy struct {
	lb          *balancer.LoadBalancer
	pool        map[string]*httputil.ReverseProxy // host:handler
	servicepath string                            // rewrite
	port        string                            // rewrite
	log         *logger.Logger
}

func (lbrp *LoadBalancedReverseProxy) next() (*httputil.ReverseProxy, error) {
	var next, err = lbrp.lb.Next()
	if err == balancer.ErrNoHostAvailable {
		return nil, err
	}
	if _, ok := lbrp.pool[next]; !ok {
		host := fmt.Sprintf("%s%s", next, lbrp.port)
		lbrp.pool[next] = &httputil.ReverseProxy{
			// see link to understand usage of rewrite
			// https://www.ory.sh/hop-by-hop-header-vulnerability-go-standard-library-reverse-proxy/
			Rewrite: func(pr *httputil.ProxyRequest) {
				pr.SetXForwarded()

				pr.Out.URL.Scheme = "https"
				pr.Out.URL.Host = host
				pr.Out.URL.Path = strings.TrimPrefix(pr.In.URL.Path, lbrp.servicepath)
				pr.Out.URL.RawPath = strings.TrimPrefix(pr.In.URL.RawPath, lbrp.servicepath)

				lbrp.log.Printf("forwarding request: (%s %s %s) => (%s %s %s)\n",
					pr.In.Method, pr.In.Host, pr.In.URL.String(),
					pr.Out.Method, pr.Out.Host, pr.Out.URL.String(),
				)
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				log.Printf("proxy error: %v, method=%s, url=%s, remoteAddr=%s\n", err, r.Method, r.URL.String(), r.RemoteAddr)
				http.Error(w, "Bad Gateway", http.StatusBadGateway)
			},
		}
	}
	return lbrp.pool[next], nil
}

func (lbrp LoadBalancedReverseProxy) Handler(w http.ResponseWriter, r *http.Request) {
	forwarder, err := lbrp.next()
	if err != nil {
		http.Error(w, "Service you want to access is not available at the moment. Please try again later.", http.StatusBadGateway)
		return
	}
	forwarder.ServeHTTP(w, r)
}

func New(sd *discovery.ServiceDiscovery, service models.Service, port, servicepath string) (*LoadBalancedReverseProxy, error) {
	var lb = balancer.New(sd, service)
	if _, err := lb.Next(); err == balancer.ErrNoHostAvailable {
		return nil, err
	}
	return &LoadBalancedReverseProxy{
		log:         logger.NewLogger("reverse proxy"),
		pool:        map[string]*httputil.ReverseProxy{},
		lb:          lb,
		servicepath: servicepath,
		port:        port,
	}, nil
}
