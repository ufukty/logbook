package captured

import (
	"cmp"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

type ResponseWriter struct {
	wrapped    http.ResponseWriter
	statusCode int
	written    bool
	size       uint
}

func New(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{wrapped: w}
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.wrapped.Header()
}

func (rw *ResponseWriter) Write(n []byte) (int, error) {
	rw.written = true
	rw.size += uint(len(n))
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

func bytes(i uint) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	j := 0
	for k := i; k >= 1024 && j+1 < len(units); k /= 1024 {
		j++
	}
	d := float64(i) / math.Pow(1024, float64(j))
	if d != float64(int(d)) {
		return fmt.Sprintf("%.1f%s", d, units[j])
	}
	return fmt.Sprintf("%d%s", int(d), units[j])
}

func (rw ResponseWriter) SizeRepresentation() string {
	return bytes(rw.size)
}
