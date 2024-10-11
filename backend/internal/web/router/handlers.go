package router

import (
	"fmt"
	"logbook/internal/logger"
	"net/http"
)

func summarize(r *http.Request) string {
	return fmt.Sprintf("(\033[34m%s\033[0m, \033[35m%s\033[0m) \033[31m%s\033[0m \033[32m%s\033[0m \033[33m%s\033[0m", r.Host, r.RemoteAddr, r.Proto, r.Method, r.URL.Path)
}

func dumpRequestBuilder(l *logger.Logger) func(r *http.Request) {
	return func(r *http.Request) {
		l.Println(summarize(r))
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
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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
