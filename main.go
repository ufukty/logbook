package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	db "logbook/main/database"

	"github.com/gorilla/mux"
)

func register_endpoints(r **mux.Router) {
	for _, endpoint := range endpoints() {
		(*r).
			HandleFunc(endpoint.Path, endpoint.Handler).
			Methods(endpoint.Method)
	}
}

func main() {
	// Make sure log package uses UTC
	log.SetFlags(log.LstdFlags | log.LUTC)

	// Accept argument from terminal for server configuration
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Initialize database connection pool
	db.Init("postgres://postgres:password@localhost:5432/testdatabase") // os.Getenv("DATABASE_URL")
	defer db.Close()

	// document.Init()

	//
	r := mux.NewRouter()

	//
	register_endpoints(&r)

	// r.Use(middleware.MWAuthorization)

	srv := &http.Server{
		Addr: ":8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
