package receptionist

import "net/http"

type Response struct {
	http.ResponseWriter
	Status int
}

func (r *Response) WriteHeader(statusCode int) {
	r.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
