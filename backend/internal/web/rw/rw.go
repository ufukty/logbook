package rw

import (
	"cmp"
	"net/http"
)

// A [http.ResponseWriter] that keeps status code accessible and supports [http.ResponseController].
type responseWriter struct {
	wrapped    http.ResponseWriter
	statusCode int
	written    bool
}

func New(rw http.ResponseWriter) http.ResponseWriter {
	return &responseWriter{wrapped: rw}
}

func (rw *responseWriter) Header() http.Header {
	return rw.wrapped.Header()
}

func (rw *responseWriter) Write(n []byte) (int, error) {
	rw.written = true
	return rw.wrapped.Write(n)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.wrapped.WriteHeader(statusCode)
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.wrapped
}

func (rw responseWriter) Status() (int, bool) {
	return cmp.Or(rw.statusCode, int(http.StatusOK)), rw.statusCode > 0 || rw.written
}
