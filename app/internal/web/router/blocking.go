package router

import (
	"fmt"
	"logbook/config/api"
	config "logbook/config/deployment"
	"logbook/internal/web/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"golang.org/x/exp/maps"
)

func StartServer(baseURL string, tls bool, cfg config.RouterParameters, handlers map[api.Endpoint]http.HandlerFunc) {
	l := logger.NewLogger("Router")
	r := mux.NewRouter()

	{
		l.Println("registering routes in the order:")
		r := r.UseEncodedPath()
		for _, ep := range sortEndpoints(maps.Keys(handlers)) {
			str := fmt.Sprintf("%s %s", ep.Method, ep.Path)
			handler := handlers[ep]
			// r.HandleFunc(str, handler)
			r.HandleFunc(string(ep.Path), handler).Methods(ep.Method)
			l.Printf("%q -> %p\n", str, handler)
		}
		r.HandleFunc("/ping", pongBuilder(l))
		r.HandleFunc("/", lastMatchBuilder(l))
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(cfg.RequestTimeout))
	r.Use(middleware.Logger)
	// r.Use(middleware.MWAuthorization)
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middleware.Recoverer)

	server := &http.Server{
		Addr: baseURL,
		// Set timeouts against Slowloris attacks
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	if tls {
		l.Printf("calling ListenAndServeTLS on '%s'\n", baseURL)
		if err := server.ListenAndServeTLS(publicCertPath, privateCertPath); err != nil {
			l.Println(fmt.Errorf("http.Server returned an error from ListenAndServeTLS call: %w", err))
		}
	} else {
		l.Printf("calling ListenAndServe on '%s'\n", baseURL)
		if err := server.ListenAndServe(); err != nil {
			l.Println(fmt.Errorf("http.Server returned an error from ListendAndServe call: %w", err))
		}
	}

}
