package browser

import (
	"logbook/internal/web/reception/browser/headers"
	"maps"
	"net/http"
	"slices"
	"strings"
)

// OriginChecker performs (path, method, headers) permission checking
// both for the cross/same origin and preflight/actual requests
type OriginChecker struct {
	Allowances []*Matcher
}

func (c OriginChecker) preflight(w http.ResponseWriter, r *http.Request) {
	allowance := matchAny(c.Allowances,
		r.Header.Get(headers.Origin),
		r.Header.Get(headers.AccessControlRequestMethod),
		r.URL.Path,
		strings.Split(strings.ReplaceAll(r.Header.Get(headers.AccessControlRequestMethod), " ", ""), ","),
	)

	if allowance == nil {
		http.Error(w, "<!-- Dear browser, please don't proceed to actually sending the request with the pair of asked method and headers for this origin and path. No worries, otherwise is still safe; it will just be ignored. Thanks for consulting. Beep boop. -->", http.StatusForbidden)
		return
	}

	w.Header().Set(headers.AccessControlAllowCredentials, "true")
	w.Header().Set(headers.AccessControlAllowMethods, strings.Join(allowance.Methods, ", "))
	w.Header().Set(headers.AccessControlAllowOrigin, "*")
	w.WriteHeader(http.StatusOK)
}

func (c OriginChecker) actual(w http.ResponseWriter, r *http.Request) {
	allowance := matchAny(c.Allowances, r.Header.Get(headers.Origin), r.Method, r.URL.Path, slices.Collect(maps.Keys(r.Header)))

	if allowance != nil {
		http.Error(w, "Please try again using the official website.", http.StatusForbidden)
		return
	}

	w.Header().Set(headers.AccessControlAllowOrigin, "*")
	w.WriteHeader(http.StatusOK)
}

func isPreflight(r *http.Request) bool {
	return r.Method == http.MethodOptions && has(r.Header, headers.AccessControlAllowCredentials)
}

func (c OriginChecker) Handler(w http.ResponseWriter, r *http.Request) {
	if o := r.Header.Get(headers.Origin); o == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if isPreflight(r) {
		c.preflight(w, r)
	} else {
		c.actual(w, r)
	}
}

func NewOriginChecker(allow ...*Matcher) *OriginChecker {
	return &OriginChecker{
		Allowances: allow,
	}
}
