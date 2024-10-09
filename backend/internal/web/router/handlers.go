package router

import (
	"context"
	"fmt"
	"log"
	"logbook/internal/logger"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func dumpRequestBuilder(l *logger.Logger) func(r *http.Request) {
	return func(r *http.Request) {
		var dump, err = httputil.DumpRequest(r, false)
		if err != nil {
			log.Println(errors.Wrap(err, "dumping request"))
		}
		l.Printf("%q", strings.ReplaceAll(strings.ReplaceAll(string(dump), "\r\n", " || "), "\n", " | "))
	}
}

func notFoundBuilder(l *logger.Logger) func(w http.ResponseWriter, r *http.Request) {
	dumpRequest := dumpRequestBuilder(l)
	return func(w http.ResponseWriter, r *http.Request) {
		dumpRequest(r)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func lastMatchBuilder(l *logger.Logger) func(w http.ResponseWriter, r *http.Request) {
	dumpRequest := dumpRequestBuilder(l)
	return func(w http.ResponseWriter, r *http.Request) {
		dumpRequest(r)
		if r.URL.Path == "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
	}
}

func pongBuilder(l *logger.Logger) func(w http.ResponseWriter, r *http.Request) {
	dumpRequest := dumpRequestBuilder(l)
	return func(w http.ResponseWriter, r *http.Request) {
		dumpRequest(r)
		fmt.Fprintln(w, "pong")
	}
}

func applyMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement request ID middleware
		next.ServeHTTP(w, r)
	})
}

func withTimeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func withLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement logger middleware
		next.ServeHTTP(w, r)
	})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement CORS middleware
		next.ServeHTTP(w, r)
	})
}

func withRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement recoverer middleware
		next.ServeHTTP(w, r)
	})
}
