// design justification:
// traditional middleware implementation (with clousures)
// doesn't support typed context, as signatures are unchangeble (w, r)
// and the flow of execution is not clear

package pipelines

import (
	"fmt"
	"logbook/internal/logger"
	"logbook/models/columns"
	"net/http"
	"runtime/debug"
)

type RequestId string

type Continuity string

const (
	Continue    = Continuity("continue")
	Error       = Continuity("error")
	EarlyReturn = Continuity("early-return")
)

type HandlerFunc[StorageType any] func(w http.ResponseWriter, r *http.Request, id RequestId, store *StorageType) Continuity

type Pipe[StorageType any] interface {
	Name() string
	Handle(w http.ResponseWriter, r *http.Request, id RequestId, store *StorageType) Continuity
}

type pipeline[StorageType any] struct {
	l       *logger.Logger
	handler HandlerFunc[StorageType]
	pre     []Pipe[StorageType]
	post    []Pipe[StorageType]
}

func NewPipeline[StorageType any](handler HandlerFunc[StorageType], pre, post []Pipe[StorageType], l *logger.Logger) *pipeline[StorageType] {
	return &pipeline[StorageType]{
		l:       l.Sub("Pipeline"),
		handler: handler,
		pre:     pre,
		post:    post,
	}
}

type phase string

const (
	pre  = phase("pre")
	in   = phase("in")
	post = phase("post")
)

// TODO: timeout middleware
func (p pipeline[StorageType]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := columns.NewUuidV4[RequestId]()
	if err != nil {
		p.l.Println(fmt.Errorf("generating new request id: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// durs := map[middleware[StorageType]]time.Time{}

	ph := pre
	var pipe Pipe[StorageType]

	defer func() {
		if rec := recover(); rec != nil {
			if rec == http.ErrAbortHandler { // don't recover
				panic(rec)
			}

			debug.PrintStack()
			logname := "handler"
			if ph != in {
				logname = pipe.Name()
			}
			p.l.Println(fmt.Errorf("recovered: %s: %v", logname, rec))

			if r.Header.Get("Connection") != "Upgrade" { // except websocket (?)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
	}()

	store := new(StorageType)

	if p.pre != nil {
		for _, pipe = range p.pre {
			ec := pipe.Handle(w, r, id, store)
			if ec == Continue {
				continue
			} else if ec == EarlyReturn {
				return
			} else {
				p.l.Println(fmt.Errorf("pre: %s: %w", pipe.Name(), err))
				return
			}
		}
	}

	ph = in
	ec := p.handler(w, r, id, store)
	if ec == EarlyReturn {
		return
	} else if ec == Error {
		p.l.Println(fmt.Errorf("handler: %w", err))
		return
	}
	ph = post

	if p.post != nil {
		for _, pipe = range p.post {
			ec := pipe.Handle(w, r, id, store)
			if ec == Continue {
				continue
			} else if ec == EarlyReturn {
				return
			} else {
				p.l.Println(fmt.Errorf("post: %s: %w", pipe.Name(), err))
				return
			}
		}
	}
}
