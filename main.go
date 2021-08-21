package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// type ThreadType string

// const (
// 	Daily        ThreadType = "Daily"
// 	Active                  = "Active"
// 	Paused                  = "Paused"
// 	ReadyToStart            = "ReadyToStart"
// 	Plan                    = "Plan"
// )

// type Task struct {
// 	content string
// 	task_id string
// }

// type ThreadTask struct {
// 	task_id    string
// 	task       Task
// 	depth      int
// 	degree     int // Number of all nodes below
// 	created_at string
// }

// type Thread struct {
// 	thread_type ThreadType
// 	tasks       []ThreadTask
// }

// type Document struct {
// 	threads       []Thread
// 	total_threads int
// }

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	for _, endpoint := range endpoints() {
		r.
			HandleFunc(endpoint.Path, endpoint.Handler).
			Methods(endpoint.Method)
	}

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
