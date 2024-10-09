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

	"github.com/go-chi/chi/v5/middleware"
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

func StartServer(params ServerParameters, endpointRegisterer func(r *http.ServeMux), l *logger.Logger) {
	tls := params.TlsKey != "" && params.TlsCrt != ""
	l.Sub("Router: ")

	mux := http.NewServeMux()
	endpointRegisterer(mux)
	mux.HandleFunc("/ping", pongBuilder(l))
	mux.HandleFunc("/", lastMatchBuilder(l))

	handler := applyMiddleware(mux,
		middleware.RequestID,
		middleware.Timeout(time.Duration(params.Router.RequestTimeout)),
		middleware.Logger,
		// middleware.MWAuthorization,
		// mux.CORSMethodMiddleware(r),
		middleware.Recoverer,
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", params.Port))
	if err != nil {
		l.Printf("net.Listen: %v\n", err)
		return
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	l.Printf("listening on port: %d\n", port)

	if params.Sidecar != nil {
		err := params.Sidecar.SetInstanceDetails(params.Service, models.Instance{
			Tls:     tls,
			Address: params.Address,
			Port:    port,
		})
		if err != nil {
			l.Printf("Sidecar.SetInstanceDetails: %v\n", err)
		}
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", params.Port),
		WriteTimeout: time.Duration(params.Router.WriteTimeout),
		ReadTimeout:  time.Duration(params.Router.ReadTimeout),
		IdleTimeout:  time.Duration(params.Router.IdleTimeout),
		Handler:      handler,
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
	signal.Notify(sigInterruptChannel, os.Interrupt)
	<-sigInterruptChannel

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(params.Router.GracePeriod))
	defer cancel()

	log.Printf("Shutting down server, grace period is '%s'\n", time.Duration(params.Router.GracePeriod).String())
	if err := server.Shutdown(ctx); err != nil {
		l.Printf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server is down")
}
