package browser

import (
	"maps"
	"net/http"
	"slices"
	"strings"
)

type allowance struct {
	methods map[string]any
	headers map[string]any
}

func lookup(ss []string, canonicalize func(string) string) map[string]any {
	us := make(map[string]any, len(ss))
	for _, s := range ss {
		us[canonicalize(s)] = nil
	}
	return us
}

func has[K comparable, V any](m map[K]V, k K) bool {
	_, ok := m[k]
	return ok
}

func contains(asked string, allowed map[string]any, canonicalize func(string) string) bool {
	return has(allowed, canonicalize(asked))
}

func containsAll(asked []string, allowed map[string]any, canonicalize func(string) string) bool {
	for _, a := range asked {
		if !contains(a, allowed, canonicalize) {
			return false
		}
	}
	return true
}

type StringMatcher interface {
	MatchString(s string) bool
}

// Matcher uses custom matcher for origin and path;
// lowercase character matching for methods and headers
type Matcher struct {
	origin    StringMatcher
	path      StringMatcher
	allowance *allowance
}

func NewMatcher(origin, path StringMatcher, allowedmethods, allowedheaders []string) *Matcher {
	return &Matcher{
		origin: origin,
		path:   path,
		allowance: &allowance{
			methods: lookup(allowedmethods, strings.ToLower),
			headers: lookup(allowedheaders, http.CanonicalHeaderKey),
		},
	}
}

type Scope struct {
	Methods []string
	Headers []string
}

func (m Matcher) Match(origin, method, path string, headers []string) *Scope {
	if origin == "" || path == "" || method == "" {
		return nil
	}
	if !m.origin.MatchString(origin) || !m.path.MatchString(path) {
		return nil
	}
	if !contains(method, m.allowance.methods, strings.ToLower) {
		return nil
	}
	if !containsAll(headers, m.allowance.headers, http.CanonicalHeaderKey) {
		return nil
	}
	return &Scope{
		Methods: slices.Collect(maps.Keys(m.allowance.methods)),
		Headers: slices.Collect(maps.Keys(m.allowance.headers)),
	}
}

func matchAny(ms []*Matcher, origin, method, path string, headers []string) *Scope {
	for _, m := range ms {
		if a := m.Match(origin, method, path, headers); a != nil {
			return a
		}
	}
	return nil
}
