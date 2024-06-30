package router

import (
	"context"
	"fmt"
	"log"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/web/logger"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"golang.org/x/exp/maps"
)

type ServerParameters struct {
	Router  deployment.Router
	BaseUrl string
	TlsCrt  string
	TlsKey  string
}

func StartServer(params ServerParameters, endpointRegisterer func(r *mux.Router)) {
	l := logger.NewLogger("Router")

	r := mux.NewRouter()
	endpointRegisterer(r)
	r.HandleFunc("/ping", pongBuilder(l))
	r.PathPrefix("/").HandlerFunc(lastMatchBuilder(l))

	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(time.Duration(params.Router.RequestTimeout)))
	r.Use(middleware.Logger)
	// r.Use(middleware.MWAuthorization)
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middleware.Recoverer)

	server := &http.Server{
		Addr: params.BaseUrl,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Duration(params.Router.WriteTimeout),
		ReadTimeout:  time.Duration(params.Router.ReadTimeout),
		IdleTimeout:  time.Duration(params.Router.IdleTimeout),
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if params.TlsKey != "" && params.TlsCrt != "" {
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
	}()

	sigInterruptChannel := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(sigInterruptChannel, os.Interrupt)
	// Block until we receive our signal.
	<-sigInterruptChannel

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(params.Router.GracePeriod))
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	log.Printf("Sending shutdown signal to one of the servers, grace period is '%s'\n", time.Duration(params.Router.GracePeriod).String())
	go server.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	<-ctx.Done()
	log.Println("Server is down")
}

func StartServerWithEndpoints(params ServerParameters, handlers map[api.Endpoint]http.HandlerFunc) {
	l := logger.NewLogger("Router")
	r := mux.NewRouter()

	{
		l.Println("registering routes in the order:")
		r := r.UseEncodedPath()
		for _, ep := range sortEndpoints(maps.Keys(handlers)) {
			str := fmt.Sprintf("%s %s", ep.GetMethod(), ep.GetPath())
			handler := handlers[ep]
			// r.HandleFunc(str, handler)
			r.HandleFunc(string(ep.GetPath()), handler).Methods(ep.GetMethod())
			l.Printf("%q -> %p\n", str, handler)
		}
		r.HandleFunc("/ping", pongBuilder(l))
		r.HandleFunc("/", lastMatchBuilder(l))
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(time.Duration(params.Router.RequestTimeout)))
	r.Use(middleware.Logger)
	// r.Use(middleware.MWAuthorization)
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middleware.Recoverer)

	server := &http.Server{
		Addr: params.BaseUrl,
		// Set timeouts against Slowloris attacks
		WriteTimeout: time.Duration(params.Router.WriteTimeout),
		ReadTimeout:  time.Duration(params.Router.ReadTimeout),
		IdleTimeout:  time.Duration(params.Router.IdleTimeout),
		Handler:      r,
	}

	go func() {
		if params.TlsKey != "" && params.TlsCrt != "" {
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
	}()

	sigInterruptChannel := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(sigInterruptChannel, os.Interrupt)
	// Block until we receive our signal.
	<-sigInterruptChannel

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(params.Router.GracePeriod))
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	log.Printf("Sending shutdown signal to one of the servers, grace period is '%s'\n", time.Duration(params.Router.GracePeriod).String())
	go server.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	<-ctx.Done()
	log.Println("Server is down")
}
