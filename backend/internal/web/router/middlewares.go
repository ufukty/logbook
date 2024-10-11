package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

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
