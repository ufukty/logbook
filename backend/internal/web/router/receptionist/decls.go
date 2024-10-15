package receptionist

import (
	"fmt"
	"net/http"
)

type RequestId string

const ZeroRequestId = RequestId("00000000-0000-0000-0000-000000000000")

var ErrEarlyReturn = fmt.Errorf("no error") // return early without logging an error

// Basically: [http.HandlerFunc] with additions
type HandlerFunc[StorageType any] func(id RequestId, store *StorageType, w http.ResponseWriter, r *http.Request) error

type Response struct {
	http.ResponseWriter
	Status int
}

func (r *Response) WriteHeader(statusCode int) {
	r.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)

}
