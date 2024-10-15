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
	"logbook/internal/web/router/registration/receptionist/decls"
	"logbook/models/columns"
	"net/http"
	"runtime/debug"
	"time"
)

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
type receptionist struct {
	l        *logger.Logger
	handlers []decls.HandlerFunc
	params   Params
}

func New(params Params, l *logger.Logger, handlers ...decls.HandlerFunc) *receptionist {
	return &receptionist{
		l:        l.Sub("Pipeline"),
		handlers: handlers,
		params:   params,
	}
}

// TODO: log
// TODO: not found
// DONE: recover
// DONE: timeout
// DONE: handlers
func (recp receptionist) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ww := &response{ResponseWriter: w}

	id, err := columns.NewUuidV4[decls.RequestId]()
	if err != nil {
		recp.l.Println(fmt.Errorf("generating new request id: %w", err))
		http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	t := time.Now()

	recp.l.Printf("accepted %s: %s\n", lastsix(id), summarize(r))
	defer func() {
		recp.l.Printf("served   %s: %s\n", lastsix(id), summarizeW(ww, t))
	}()

	ctx, cancel := context.WithTimeout(r.Context(), recp.params.Timeout)
	defer func() {
		cancel()
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(ww, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		}
	}()
	r = r.WithContext(ctx)

	var handler decls.HandlerFunc
	defer func() {
		if rec := recover(); rec != nil {
			if rec == http.ErrAbortHandler { // don't recover
				panic(rec)
			}
			debug.PrintStack()
			recp.l.Println(fmt.Errorf("recovered: %s: %v", funcname(handler), rec))
			if r.Header.Get("Connection") != "Upgrade" { // except websocket (?)
				http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
	}()

	select {
	case <-r.Context().Done(): // handle timeout
		return

	default:
		store := &decls.Store{}
		for _, handler = range recp.handlers {
			err := handler(id, store, ww, r)
			if err != nil {
				if err != decls.ErrEarlyReturn {
					recp.l.Println(fmt.Errorf("handler %s: %w", funcname(handler), err))
				}
				return
			}
		}
	}

}
