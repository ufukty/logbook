// design justification:
// traditional middleware implementation (with clousures)
// doesn't support typed context, as signatures are unchangeble (w, r)
// and the flow of execution is not clear

// the timeout and recovery logic are based on chi's middlewares

package reception

import (
	"context"
	"fmt"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/models/columns"
	"net/http"
	"runtime/debug"
	"time"
)

type RequestId string

const ZeroRequestId = RequestId("00000000-0000-0000-0000-000000000000")

type receptionist struct {
	l       *logger.Logger
	handler http.Handler
	deplcfg *deployment.Config
}

func newReceptionist(deplcfg *deployment.Config, l *logger.Logger, handler http.Handler) *receptionist {
	return &receptionist{
		l:       l.Sub("receptionist"),
		handler: handler,
		deplcfg: deplcfg,
	}
}

// DONE: logging
// DONE: recover
// DONE: timeout
func (recp receptionist) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ww := &response{ResponseWriter: w}

	id, err := columns.NewUuidV4[RequestId]()
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

	ctx, cancel := context.WithTimeout(r.Context(), recp.deplcfg.Reception.RequestTimeout)
	defer func() {
		cancel()
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(ww, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		}
	}()
	r = r.WithContext(ctx)

	defer func() {
		if rec := recover(); rec != nil {
			if rec == http.ErrAbortHandler { // don't recover
				panic(rec)
			}
			debug.PrintStack()
			recp.l.Println(fmt.Errorf("recovered: %s: %v", funcname(recp.handler), rec))
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
		recp.handler.ServeHTTP(ww, r)
	}

}
