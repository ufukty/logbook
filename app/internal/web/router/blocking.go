package router

import (
	"context"
	"fmt"
	"log"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/web/discoveryctl"
	"logbook/internal/web/logger"
	"logbook/models"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"golang.org/x/exp/maps"
)

type ServerParameters struct {
	Discovery *discoveryctl.Client
	Service   models.Service
	Address   string
	Router    deployment.Router
	Port      int
	TlsCrt    string
	TlsKey    string
}

func StartServer(params ServerParameters, endpointRegisterer func(r *mux.Router)) {
	tls := params.TlsKey != "" && params.TlsCrt != ""
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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", params.Port))
	if err != nil {
		l.Println(fmt.Errorf("net.Listen: %w", err))
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	l.Printf("listening the port: %d\n", port)

	if params.Discovery != nil {
		err := params.Discovery.SetInstanceDetails(params.Service, models.Instance{
			Tls:     tls,
			Address: params.Address,
			Port:    port,
		})
		if err != nil {
			l.Println(fmt.Errorf("params.Discovery.SetInstanceDetails: %w", err))
		}
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", params.Port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Duration(params.Router.WriteTimeout),
		ReadTimeout:  time.Duration(params.Router.ReadTimeout),
		IdleTimeout:  time.Duration(params.Router.IdleTimeout),
		Handler:      r,
	}

	go func() {
		if tls {
			if err := server.ServeTLS(listener, params.TlsCrt, params.TlsKey); err != nil {
				l.Println(fmt.Errorf("server.ServeTLS: %w", err))
			}
		} else {
			if err := server.Serve(listener); err != nil {
				l.Println(fmt.Errorf("server.Serve: %w", err))
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
	l := logger.NewLogger("StartServerWithEndpoints")
	StartServer(params, func(r *mux.Router) {
		l.Println("registering routes in the order:")
		r = r.UseEncodedPath()
		for _, ep := range sortEndpoints(maps.Keys(handlers)) {
			str := fmt.Sprintf("%s %s", ep.GetMethod(), ep.GetPath())
			handler := handlers[ep]
			// r.HandleFunc(str, handler)
			r.HandleFunc(string(ep.GetPath()), handler).Methods(ep.GetMethod())
			l.Printf("%q -> %p\n", str, handler)
		}
	})
}
