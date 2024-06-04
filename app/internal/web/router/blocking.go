package router

import (
	"fmt"
	"log"
	"logbook/config/api"
	"logbook/internal/web/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"golang.org/x/exp/maps"
)

type ServerParameters struct {
	BaseUrl        string
	Tls            bool
	TlsCrt         string
	TlsKey         string
	RequestTimeout time.Duration
}

func StartServer(params ServerParameters, endpointRegisterer func(r *mux.Router)) {
	l := logger.NewLogger("Router")

	r := mux.NewRouter()
	endpointRegisterer(r)
	r.HandleFunc("/ping", pongBuilder(l))
	r.PathPrefix("/").HandlerFunc(lastMatchBuilder(l))

	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(params.RequestTimeout))
	r.Use(middleware.Logger)
	// r.Use(middleware.MWAuthorization)
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middleware.Recoverer)

	server := &http.Server{
		Addr: params.BaseUrl,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	if params.Tls {
		l.Printf("calling ListenAndServeTLS on %q\n", params.BaseUrl)
		if err := server.ListenAndServeTLS(params.TlsCrt, params.TlsKey); err != nil {
			l.Println(fmt.Errorf("http.Server returned an error from ListenAndServeTLS call: %w", err))
		}
	} else {
		l.Printf("calling ListenAndServe on %q\n", params.BaseUrl)
		if err := server.ListenAndServe(); err != nil {
			l.Println(fmt.Errorf("http.Server returned an error from ListendAndServe call: %w", err))
		}
	}
}

func StartServerWithEndpoints(params ServerParameters, handlers map[api.Endpoint]http.HandlerFunc) {
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
	r.Use(middleware.Timeout(params.RequestTimeout))
	r.Use(middleware.Logger)
	// r.Use(middleware.MWAuthorization)
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middleware.Recoverer)

	server := &http.Server{
		Addr: params.BaseUrl,
		// Set timeouts against Slowloris attacks
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	if params.Tls {
		l.Printf("calling ListenAndServeTLS on %q\n", params.BaseUrl)
		if err := server.ListenAndServeTLS(params.TlsCrt, params.TlsKey); err != nil {
			l.Println(fmt.Errorf("http.Server returned an error from ListenAndServeTLS call: %w", err))
		}
	} else {
		l.Printf("calling ListenAndServe on %q\n", params.BaseUrl)
		if err := server.ListenAndServe(); err != nil {
			l.Println(fmt.Errorf("http.Server returned an error from ListendAndServe call: %w", err))
		}
	}
}
