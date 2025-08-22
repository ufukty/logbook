package browser

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
)

func s[T any](ss ...T) []T {
	return ss
}

func TestMatcherMatch(t *testing.T) {
	// cases belong scenarios
	type (
		tc struct {
			Origin  string
			Path    string
			Method  string
			Headers []string
		}
		ts struct {
			Matcher  *Matcher // shared among cases
			Positive map[string]tc
			Negative map[string]tc
		}
	)
	tsc := map[string]ts{
		"hardcoded origin and path": {
			Matcher: NewMatcher(regexp.MustCompile("localhost:3000"), regexp.MustCompile("/user"), s(http.MethodGet), []string{"Content-Type"}),
			Positive: map[string]tc{
				"": {"localhost:3000", "/user", http.MethodGet, []string{"Content-Type"}},
			},
			Negative: map[string]tc{
				"cross origin": {"localhost:8080", "/user", http.MethodGet, []string{}},
				"cross path":   {"localhost:3000", "/account", http.MethodGet, []string{}},
			},
		},
		"wildcard origin and path": {
			Matcher: NewMatcher(regexp.MustCompile(".*"), regexp.MustCompile(".*"), s(http.MethodGet), s("Content-Type")),
			Positive: map[string]tc{
				"member": {"localhost", "/user", http.MethodGet, s("Content-Type")},
			},
			Negative: map[string]tc{
				"empty method":     {"localhost", "/user", "", s("Content-Type")},
				"unallowed method": {"localhost", "/user", http.MethodPost, s("Content-Type")},
				"unallowed header": {"localhost", "/user", http.MethodGet, s("Cookie")},
			},
		},
		"no headers": {
			Matcher: NewMatcher(regexp.MustCompile("localhost:3000"), regexp.MustCompile("/user"), s(http.MethodGet), []string{}),
			Positive: map[string]tc{
				"": {"localhost:3000", "/user", http.MethodGet, []string{}},
			},
			Negative: map[string]tc{
				"cross origin": {"localhost:8080", "/user", http.MethodGet, []string{}},
				"cross path":   {"localhost:3000", "/account", http.MethodGet, []string{}},
			},
		},
	}

	for tsn, ts := range tsc {
		for tcn, tc := range ts.Positive {
			t.Run(fmt.Sprintf("%q should ALLOW %q", tsn, tcn), func(t *testing.T) {
				if a := ts.Matcher.Match(tc.Origin, tc.Path, tc.Method, tc.Headers); a == nil {
					t.Errorf("expected match")
				}
			})
		}

		for tcn, tc := range ts.Negative {
			t.Run(fmt.Sprintf("%q should UNALLOW %q", tsn, tcn), func(t *testing.T) {
				if a := ts.Matcher.Match(tc.Origin, tc.Path, tc.Method, tc.Headers); a != nil {
					t.Errorf("unexpected match")
				}
			})
		}
	}
}
