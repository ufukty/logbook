package decls

import (
	"fmt"
	"logbook/models/columns"
	"net/http"
)

type RequestId string

const ZeroRequestId = RequestId("00000000-0000-0000-0000-000000000000")

var ErrEarlyReturn = fmt.Errorf("no error") // return early without logging an error

// Basically: [http.HandlerFunc] with additions for type safety
type HandlerFunc func(id RequestId, store *Store, w http.ResponseWriter, r *http.Request) error

type Store struct {
	User columns.UserId
}
