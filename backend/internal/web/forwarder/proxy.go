package forwarder

import (
	"fmt"
	"logbook/internal/logger"
	"logbook/internal/web/balancer"
	"logbook/models"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

type LoadBalancedReverseProxy struct {
	lb          *balancer.LoadBalancer
	pool        map[*models.Instance]*httputil.ReverseProxy
	mu          sync.RWMutex
	servicepath string // rewrite
	log         *logger.Logger
}

func (lbrp *LoadBalancedReverseProxy) next() (*httputil.ReverseProxy, error) {
	var next, err = lbrp.lb.Next()
	if err == balancer.ErrNoHostAvailable {
		return nil, err
	}
	lbrp.mu.RLock()
	nextrp, ok := lbrp.pool[next]
	lbrp.mu.RUnlock()
	if !ok {
		host := fmt.Sprintf("%s:%d", next.Address, next.Port)
		nextrp = &httputil.ReverseProxy{
			// see link to understand usage of rewrite
			// https://www.ory.sh/hop-by-hop-header-vulnerability-go-standard-library-reverse-proxy/
			Rewrite: func(pr *httputil.ProxyRequest) {
				pr.SetXForwarded()

				pr.Out.URL.Scheme = "https"
				pr.Out.URL.Host = host
				pr.Out.URL.Path = strings.TrimPrefix(pr.In.URL.Path, lbrp.servicepath)
				pr.Out.URL.RawPath = strings.TrimPrefix(pr.In.URL.RawPath, lbrp.servicepath)

				pr.Out.Host = pr.In.Host

				// lbrp.log.Printf("forwarding request: (%s %s %s) => (%s %s %s)\n",
				// 	pr.In.Method, pr.In.Host, pr.In.URL.String(),
				// 	pr.Out.Method, pr.Out.Host, pr.Out.URL.String(),
				// )
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				lbrp.log.Printf("proxy error: %v, method=%s, url=%s, remoteAddr=%s\n", err, r.Method, r.URL.String(), r.RemoteAddr)
				http.Error(w, "Bad Gateway", http.StatusBadGateway)
			},
		}
		lbrp.mu.Lock()
		lbrp.pool[next] = nextrp
		lbrp.mu.Unlock()
	}
	return nextrp, nil
}

func (lbrp *LoadBalancedReverseProxy) Handler(w http.ResponseWriter, r *http.Request) {
	forwarder, err := lbrp.next()
	if err != nil {
		lbrp.log.Println(fmt.Errorf("lbrp.next: %w", err))
		http.Error(w, "Service you want to access is not available at the moment. Please try again later.", http.StatusBadGateway)
		return
	}
	forwarder.ServeHTTP(w, r)
}

func New(is balancer.InstanceSource, service models.Service, servicepath string) *LoadBalancedReverseProxy {
	return &LoadBalancedReverseProxy{
		log:         logger.NewLogger("LoadBalancedReverseProxy"),
		pool:        map[*models.Instance]*httputil.ReverseProxy{},
		lb:          balancer.New(is),
		servicepath: servicepath,
	}
}
