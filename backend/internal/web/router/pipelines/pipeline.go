// design justification:
// traditional middleware implementation (with clousures)
// doesn't support typed context, as signatures are unchangeble (w, r)
// and the flow of execution is not clear

package pipelines

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

var ErrSilent = fmt.Errorf("no error") // return early without logging an error

type HandlerFunc[StorageType any] func(id RequestId, store *StorageType, w http.ResponseWriter, r *http.Request) error

type PipelineParams struct {
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
type pipeline[StorageType any] struct {
	l        *logger.Logger
	handlers []HandlerFunc[StorageType]
	params   PipelineParams
}

func NewPipeline[T any](handlers []HandlerFunc[T], params PipelineParams, l *logger.Logger) *pipeline[T] {
	return &pipeline[T]{
		l:        l.Sub("Pipeline"),
		handlers: handlers,
		params:   params,
	}
}

func (p pipeline[StorageType]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := columns.NewUuidV4[RequestId]()
	if err != nil {
		p.l.Println(fmt.Errorf("generating new request id: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	p.l.Printf("accepted %s: %s\n", lastsix(id), summarize(r))
	defer func() {
		p.l.Printf("served %s: %s\n", lastsix(id), summarize(r))
	}()

	ctx, cancel := context.WithTimeout(r.Context(), p.params.Timeout)
	defer func() {
		cancel()
		if ctx.Err() == context.DeadlineExceeded {
			w.WriteHeader(http.StatusGatewayTimeout)
		}
	}()
	r = r.WithContext(ctx)

	defer func() {
		if rec := recover(); rec != nil {
			if rec == http.ErrAbortHandler { // don't recover
				panic(rec)
			}

			debug.PrintStack()
			p.l.Println(fmt.Errorf("recovered: %v", rec))

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
		for i, pipe := range p.handlers {
			err := pipe(id, store, w, r)
			if err != nil {
				if err != ErrSilent {
					p.l.Println(fmt.Errorf("handler %d: %w", i, err))
				}
				return
			}
		}
	}

}
