// design justification:
// traditional middleware implementation (with clousures)
// doesn't support typed context, as signatures are unchangeble (w, r)
// and the flow of execution is not clear

// the timeout and recovery logic are based on chi's middlewares

package reception

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"

	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/web/registration/reception/captured"
	"logbook/internal/web/registration/reception/summarizer"
	"logbook/models/columns"
)

type RequestId string

func (id RequestId) lastsix() string {
	return string(id)[max(0, len(string(id))-6):]
}

const ZeroRequestId = RequestId("00000000-0000-0000-0000-000000000000")

func funcname(i any) string {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Func {
		return "(Not a function)"
	}
	f := runtime.FuncForPC(reflect.ValueOf(i).Pointer())
	if f == nil {
		return "(Unknown function)"
	}
	return f.Name()
}

type receptionist struct {
	c *deployment.Config
	s *summarizer.Summarizer
	l *logger.Logger
	h http.Handler
}

func newReceptionist(c *deployment.Config, l *logger.Logger, h http.Handler) *receptionist {
	return &receptionist{
		c: c,
		s: summarizer.New(c.Environment == "local"),
		l: l.Sub("receptionist"),
		h: h,
	}
}

// DONE: logging
// DONE: recover
// DONE: timeout
func (rc receptionist) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	crw := captured.New(w)

	id, err := columns.NewUuidV4[RequestId]()
	if err != nil {
		rc.l.Println(fmt.Errorf("generating new request id: %w", err))
		http.Error(crw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	t := time.Now()

	rc.l.Printf("accepted %s: %s\n", id.lastsix(), rc.s.Pre(r))
	defer func() { rc.l.Printf("served   %s: %s\n", id.lastsix(), rc.s.Post(crw, t)) }()

	ctx, cancel := context.WithTimeout(r.Context(), rc.c.Reception.RequestTimeout)
	defer func() {
		cancel()
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(crw, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		}
	}()
	r = r.WithContext(ctx)

	defer func() {
		if rec := recover(); rec != nil {
			if rec == http.ErrAbortHandler { // don't recover
				panic(rec)
			}
			debug.PrintStack()
			rc.l.Println(fmt.Errorf("recovered: %s: %v", funcname(rc.h), rec))
			if r.Header.Get("Connection") != "Upgrade" { // except websocket (?)
				http.Error(crw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
	}()

	select {
	case <-r.Context().Done(): // handle timeout
		return

	default:
		rc.h.ServeHTTP(crw, r)
	}
}
