package browser

import (
	"maps"
	"net/http"
	"slices"
	"strings"
)

const (
	origin                        = "Origin"
	accessControlAllowCredentials = "Access-Control-Allow-Credentials"
	accessControlAllowMethods     = "Access-Control-Allow-Methods"
	accessControlAllowOrigin      = "Access-Control-Allow-Origin"
	accessControlRequestMethod    = "Access-Control-Request-Method"
)

// OriginChecker performs (path, method, headers) permission checking
// both for the cross/same origin and preflight/actual requests
type OriginChecker struct {
	Allowances []*Matcher
}

func (c OriginChecker) preflight(w http.ResponseWriter, r *http.Request) {
	allowance := matchAny(c.Allowances,
		r.Header.Get(origin),
		r.Header.Get(accessControlRequestMethod),
		r.URL.Path,
		strings.Split(strings.ReplaceAll(r.Header.Get(accessControlRequestMethod), " ", ""), ","),
	)

	if allowance == nil {
		http.Error(w, "<!-- Dear browser, please don't proceed to actually sending the request with the pair of asked method and headers for this origin and path. No worries, otherwise is still safe; it will just be ignored. Thanks for consulting. Beep boop. -->", http.StatusForbidden)
		return
	}

	w.Header().Set(accessControlAllowCredentials, "true")
	w.Header().Set(accessControlAllowMethods, strings.Join(allowance.Methods, ", "))
	w.Header().Set(accessControlAllowOrigin, "*")
	w.WriteHeader(http.StatusOK)
}

func (c OriginChecker) actual(w http.ResponseWriter, r *http.Request) {
	allowance := matchAny(c.Allowances, r.Header.Get(origin), r.Method, r.URL.Path, slices.Collect(maps.Keys(r.Header)))

	if allowance != nil {
		http.Error(w, "Please try again using the official website.", http.StatusForbidden)
		return
	}

	w.Header().Set(accessControlAllowOrigin, "*")
	w.WriteHeader(http.StatusOK)
}

func isPreflight(r *http.Request) bool {
	return r.Method == http.MethodOptions && has(r.Header, accessControlAllowCredentials)
}

func (c OriginChecker) Handler(w http.ResponseWriter, r *http.Request) {
	if o := r.Header.Get(origin); o == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if isPreflight(r) {
		c.preflight(w, r)
	} else {
		c.actual(w, r)
	}
}
