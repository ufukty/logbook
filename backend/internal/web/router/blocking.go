package router

import (
	"context"
	"fmt"
	"log"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/sidecar"
	"logbook/models"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
)

type ServerParameters struct {
	Address string
	Port    int
	Router  deployment.Router
	Service models.Service
	Sidecar *sidecar.Sidecar
	TlsCrt  string
	TlsKey  string
}

func StartServer(params ServerParameters, endpointRegisterer func(r *mux.Router), l *logger.Logger) {
	tls := params.TlsKey != "" && params.TlsCrt != ""
	l = l.Sub("Router")

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

	if params.Sidecar != nil {
		err := params.Sidecar.SetInstanceDetails(params.Service, models.Instance{
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
