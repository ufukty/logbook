package router

import (
	"fmt"
	"logbook/config"
	"logbook/internal/web/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
)

func StartServer(baseURL string, tls bool, cfg config.RouterParameters, endpointRegisterer func(r *mux.Router)) {
	l := logger.NewLogger(fmt.Sprintf("router(%s)", baseURL))

	r := mux.NewRouter()
	endpointRegisterer(r)
	r.HandleFunc("/ping", pongBuilder(l))
	r.PathPrefix("/").HandlerFunc(lastMatchBuilder(l))

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
