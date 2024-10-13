package router

import (
	"context"
	"fmt"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/sidecar"
	"logbook/models"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Registerer func(r *http.ServeMux) error

type ServerParameters struct {
	Address  string
	Port     int
	Router   deployment.Router
	ServeMux *http.ServeMux
	Service  models.Service
	Sidecar  *sidecar.Sidecar
	TlsCrt   string
	TlsKey   string
}

func StartServer(params ServerParameters, l *logger.Logger) error {
	tls := params.TlsKey != "" && params.TlsCrt != ""
	l.Sub("Router: ")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", params.Port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
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
			return fmt.Errorf("Sidecar.SetInstanceDetails: %w", err)
		}
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", params.Port),
		WriteTimeout: time.Duration(params.Router.WriteTimeout),
		ReadTimeout:  time.Duration(params.Router.ReadTimeout),
		IdleTimeout:  time.Duration(params.Router.IdleTimeout),
		Handler:      params.ServeMux,
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

	l.Printf("Shutting down server, grace period is '%s'\n", time.Duration(params.Router.GracePeriod).String())
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	l.Println("Server is down")
	return nil
}
