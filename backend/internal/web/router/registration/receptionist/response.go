package receptionist

import "net/http"

type response struct {
	http.ResponseWriter
	Status int
}

func (r *response) WriteHeader(statusCode int) {
	r.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
