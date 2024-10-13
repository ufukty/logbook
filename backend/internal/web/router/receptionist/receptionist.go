// design justification:
// traditional middleware implementation (with clousures)
// doesn't support typed context, as signatures are unchangeble (w, r)
// and the flow of execution is not clear

// the timeout and recovery logic are based on chi's middlewares

package receptionist

import (
	"context"
	"fmt"
	"logbook/internal/logger"
	"logbook/models/columns"
	"net/http"
	"runtime/debug"
	"time"
)

type RequestId string

var ErrEarlyReturn = fmt.Errorf("no error") // return early without logging an error

// Basically: [http.HandlerFunc] with additions
type HandlerFunc[StorageType any] func(id RequestId, store *StorageType, w http.ResponseWriter, r *http.Request) error

type Params struct {
	Timeout time.Duration
}

// Pipeline:
//
//   - accepts a request,
//   - generates request id,
//   - inits type-safe storage,
//   - calls handlers in order,
//   - checks timeout,
//   - recovers panic,
//   - handles logging
type receptionist[StorageType any] struct {
	l        *logger.Logger
	handlers []HandlerFunc[StorageType]
	params   Params
}

func New[T any](params Params, l *logger.Logger, handlers ...HandlerFunc[T]) *receptionist[T] {
	return &receptionist[T]{
		l:        l.Sub("Pipeline"),
		handlers: handlers,
		params:   params,
	}
}

func (recp receptionist[StorageType]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := columns.NewUuidV4[RequestId]()
	if err != nil {
		recp.l.Println(fmt.Errorf("generating new request id: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	recp.l.Printf("accepted %s: %s\n", lastsix(id), summarize(r))
	defer func() {
		recp.l.Printf("served %s: %s\n", lastsix(id), summarize(r))
	}()

	ctx, cancel := context.WithTimeout(r.Context(), recp.params.Timeout)
	defer func() {
		cancel()
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		}
	}()
	r = r.WithContext(ctx)

	var handler HandlerFunc[StorageType]
	defer func() {
		if rec := recover(); rec != nil {
			if rec == http.ErrAbortHandler { // don't recover
				panic(rec)
			}
			debug.PrintStack()
			recp.l.Println(fmt.Errorf("recovered: %s: %v", funcname(handler), rec))
			if r.Header.Get("Connection") != "Upgrade" { // except websocket (?)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
	}()

	select {
	case <-r.Context().Done(): // handle timeout
		return

	default:
		store := new(StorageType)
		for _, handler = range recp.handlers {
			err := handler(id, store, w, r)
			if err != nil {
				if err != ErrEarlyReturn {
					recp.l.Println(fmt.Errorf("handler %s: %w", funcname(handler), err))
				}
				return
			}
		}
	}

}
