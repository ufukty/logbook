package forwarder

import (
	"fmt"
	"logbook/internal/web/balancer"
	"logbook/internal/web/discovery"
	"logbook/models"

	"net/http"
	"net/http/httputil"
	"strings"
)

type reverseProxyPool struct {
	lb        *balancer.LoadBalancer
	pool      map[string]*httputil.ReverseProxy // host:handler
	generator func(string) *httputil.ReverseProxy
}

func newPool(lb *balancer.LoadBalancer, generator func(host string) *httputil.ReverseProxy) *reverseProxyPool {
	return &reverseProxyPool{
		pool:      map[string]*httputil.ReverseProxy{},
		lb:        lb,
		generator: generator,
	}
}

func (rpp *reverseProxyPool) Get() (*httputil.ReverseProxy, error) {
	var next, err = rpp.lb.Next()
	if err == balancer.ErrNoHostAvailable {
		return nil, err
	}
	if _, ok := rpp.pool[next]; !ok {
		rpp.pool[next] = rpp.generator(next)
	}
	return rpp.pool[next], nil
}

func newReverseProxyGenerator(servicepath, port string) func(host string) *httputil.ReverseProxy {
	return func(host string) *httputil.ReverseProxy {
		target := fmt.Sprintf("%s%s", host, port)
		// see link to understand usage of rewrite
		// https://www.ory.sh/hop-by-hop-header-vulnerability-go-standard-library-reverse-proxy/
		return &httputil.ReverseProxy{
			Rewrite: func(pr *httputil.ProxyRequest) {
				pr.SetXForwarded()

				pr.Out.URL.Scheme = "http"
				pr.Out.URL.Host = target
				pr.Out.URL.Path = strings.TrimPrefix(pr.In.URL.Path, servicepath)
				pr.Out.URL.RawPath = strings.TrimPrefix(pr.In.URL.RawPath, servicepath)
			},
		}
	}
}

func generateLoadBalancedProxyHandler(pool *reverseProxyPool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var forwarder, err = pool.Get()
		if err != nil {
			http.Error(w, "Service you want to access is not available at the moment. Please try again later.", http.StatusBadGateway)
			return
		}
		forwarder.ServeHTTP(w, r)
	}
}

func NewLoadBalancedProxy(sd *discovery.ServiceDiscovery, service models.Service, port, servicepath string) (http.HandlerFunc, error) {
	var lb = balancer.New(sd, service)
	if _, err := lb.Next(); err == balancer.ErrNoHostAvailable {
		return nil, err
	}
	var pool = newPool(lb, newReverseProxyGenerator(servicepath, port))
	return generateLoadBalancedProxyHandler(pool), nil
}
