package captured

import (
	"cmp"
	"net/http"
	"strconv"
)

type ResponseWriter struct {
	wrapped    http.ResponseWriter
	statusCode int
	written    bool
}

func New(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{wrapped: w}
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.wrapped.Header()
}

func (rw *ResponseWriter) Write(n []byte) (int, error) {
	rw.written = true
	return rw.wrapped.Write(n)
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.wrapped.WriteHeader(statusCode)
}

// Supports [http.ResponseController]
func (rw *ResponseWriter) Unwrap() http.ResponseWriter {
	return rw.wrapped
}

func (rw ResponseWriter) Status() (int, bool) {
	return cmp.Or(rw.statusCode, int(http.StatusOK)), rw.statusCode > 0 || rw.written
}

func (rw ResponseWriter) StatusRepresentation() string {
	if !rw.written {
		return strconv.Itoa(rw.statusCode) + "*"
	}
	return strconv.Itoa(rw.statusCode)
}
