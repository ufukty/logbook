package endpoints

import (
	"context"
	"fmt"
	"logbook/cmd/objectives/api/public/app"
	"logbook/internal/logger"
	"logbook/internal/rates"
	"logbook/models/columns"
	"net/http"
	"time"
)

type Endpoints struct {
	a *app.App
	l *logger.Logger
	r *rates.Layered
}

func New(a *app.App, l *logger.Logger) *Endpoints {
	return &Endpoints{
		a: a,
		l: l.Sub("Endpoints"),
		r: rates.NewLayered(rates.LimiterParams{}, rates.LimiterParams{}),
	}
}

type RequestId string

type RequestContext struct {
	Rid   RequestId
	Start time.Time
}

type gatekeeper struct {
}

func (e *Endpoints) Handler(w http.ResponseWriter, r *http.Request) {
	id, err := columns.NewUuidV4[RequestId]()
	if err != nil {
		e.l.Println(fmt.Errorf("generating new request: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			e.l.Println(fmt.Errorf("recovering from panic: %v", r))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}()

	ctx, _ := context.WithTimeout(r.Context(), max(50*time.Millisecond, timeout-500*time.Millisecond))

	err = e.r.Allow(ctx, uid)

	q := &RequestContext{
		Rid:   id,
		Start: time.Now(),
	}

	defer func() {
		e.r.Judge(ctx, uid, err)
	}()

	var endpoint http.HandlerFunc

	endpoint(w, r)
}

type Handler func(http.ResponseWriter, *http.Request, RequestContext)
