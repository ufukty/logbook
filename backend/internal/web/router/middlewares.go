package router

import (
	"logbook/internal/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func Timeout(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, time.Second*2, "")
}

func Logger(next http.Handler) http.Handler {
	var log = logger.New("Router/Middleware/Logger")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println()
	})
}

func SimpleHandler(r *http.ServeMux, timeout time.Duration) http.Handler {
	return applyMiddleware(r,
		middleware.RequestID,
		middleware.Timeout(timeout),
		middleware.Logger,
		// middleware.MWAuthorization,
		// mux.CORSMethodMiddleware(r),
		middleware.Recoverer,
	)
}
