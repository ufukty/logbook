package router

// import (
// 	"context"
// 	"fmt"
// 	"logbook/config/deployment"
// 	"logbook/internal/web/logger"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"slices"
// 	"time"

// 	"github.com/go-chi/chi/middleware"
// 	"github.com/gorilla/mux"
// )

// type Router struct {
// 	server *http.Server
// 	log    *logger.Logger
// 	cfg    config.RouterParameters
// }

// func NewRouter(baseURL string, cfg config.RouterParameters, endpointRegisterer func(r *mux.Router)) *Router {
// 	l := logger.NewLogger("Router")
// 	r := &Router{
// 		log: l,
// 		cfg: cfg,
// 	}

// 	m := mux.NewRouter()
// 	endpointRegisterer(m)
// 	m.HandleFunc("/ping", pongBuilder(l))
// 	m.PathPrefix("/").HandlerFunc(lastMatchBuilder(l))

// 	m.Use(middleware.RequestID)
// 	m.Use(middleware.Timeout(cfg.RequestTimeout))
// 	m.Use(middleware.Logger)
// 	// r.Use(middleware.MWAuthorization)
// 	m.Use(mux.CORSMethodMiddleware(m))
// 	m.Use(middleware.Recoverer)

// 	r.server = &http.Server{
// 		Addr: baseURL,
// 		// Set timeouts against Slowloris attacks
// 		WriteTimeout: cfg.WriteTimeout,
// 		ReadTimeout:  cfg.ReadTimeout,
// 		IdleTimeout:  cfg.IdleTimeout,
// 		Handler:      m,
// 	}

// 	return r
// }

// func (*Router) StartRouter() {}

// func (m *Router) StartNewRouter(baseURL string, cfg config.RouterParameters, endpointRegisterer func(r *mux.Router)) {
// 	m.log.Printf("Calling ListenAndServe on '%s'\n", baseURL)
// 	if err := m.server.ListenAndServe(); err != nil {
// 		m.log.Println(fmt.Errorf("http.Server returned an error from ListendAndServe call", err))
// 	}
// }

// func (m *Router) StartNewTLSRouter(baseURL string, endpointRegisterer func(r *mux.Router)) {
// 	r := mux.NewRouter()
// 	endpointRegisterer(r)
// 	// r.Use(middleware.MWAuthorization)
// 	r.Use(mux.CORSMethodMiddleware(r))

// 	server := &http.Server{
// 		Addr: baseURL,
// 		// Set timeouts against Slowloris attacks
// 		WriteTimeout: time.Second * 15,
// 		ReadTimeout:  time.Second * 15,
// 		IdleTimeout:  time.Second * 60,
// 		Handler:      r,
// 	}

// 	go func() {
// 		m.log.Printf("calling ListenAndServeTLS on '%s'\n", baseURL)
// 		if err := server.ListenAndServeTLS(publicCertPath, privateCertPath); err != nil {
// 			m.log.Println(fmt.Errorf("http.Server returned an error from ListenAndServeTLS call", err))
// 		}
// 		m.returning <- server
// 		m.returned = append(m.returned, server)
// 	}()

// 	m.servers = append(m.servers, server)
// }

// func (m *Router) gracefull() {
// 	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.GracePeriod)
// 	defer cancel()

// 	for _, server := range m.servers {
// 		if slices.Contains(m.returned, server) {
// 			continue
// 		}
// 		m.log.Printf("shuting down a server (grace period is '%s')\n", m.cfg.GracePeriod.String())
// 		go server.Shutdown(ctx)
// 	}

// 	select {
// 	case <-ctx.Done():
// 		m.log.Println("all servers are closed at the end of gracefull termination period")
// 		return
// 	case <-m.returning:
// 		m.log.Println("all servers are closed before gracefull termination period ends")
// 		if len(m.returned) == len(m.servers) {
// 			return
// 		}
// 	}
// }

// func (m *Router) Foreground() {
// 	sigInterruptChannel := make(chan os.Signal, 1)
// 	signal.Notify(sigInterruptChannel, os.Interrupt)
// 	<-sigInterruptChannel

// 	m.gracefull()
// }
