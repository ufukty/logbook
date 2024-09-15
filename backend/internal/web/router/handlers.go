package router

import (
	"fmt"
	"log"
	"logbook/internal/web/logger"
	"net/http"
	"net/http/httputil"
	"strings"

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
