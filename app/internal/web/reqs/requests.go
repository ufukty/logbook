package reqs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"logbook/internal/utilities/mapw"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

func parseRequestByStructFields(r *http.Request, bq any) error {
	var (
		bqt       = reflect.TypeOf(bq).Elem()
		bqv       = reflect.ValueOf(bq).Elem()
		bqvfields = bqv.NumField()
	)
	var (
		missingUrlParameters  = []string{}
		mismatchUrlParameters = []string{}
		missingCookies        = []string{}
		mismatchCookies       = []string{}
	)
	var (
		vars       = mux.Vars(r)
		cookies    = mapw.Mapify(r.Cookies(), func(c *http.Cookie) string { return c.Name })
		cookieType = reflect.TypeOf(&http.Cookie{})
	)
	var (
		cookievalue *http.Cookie
		urlvalue    string
		ok          bool
		metakey     string
	)

	for i := 0; i < bqvfields; i++ {
		sf := bqt.Field(i)
		fv := bqv.Field(i)

		if metakey, ok = sf.Tag.Lookup("cookie"); ok {
			if cookievalue, ok = cookies[metakey]; ok {
				if cookieType.ConvertibleTo(fv.Type()) {
					fv.Set(reflect.ValueOf(cookievalue.Value))
				} else {
					mismatchCookies = append(mismatchCookies, fmt.Sprintf("%q (%q -> %q)", metakey, "string", sf.Type.String()))
				}
			} else {
				missingCookies = append(missingCookies, metakey)
			}
		} else if metakey, ok = sf.Tag.Lookup("url"); ok {
			if urlvalue, ok = vars[metakey]; ok {
				if cookieType.ConvertibleTo(fv.Type()) {
					fv.Set(reflect.ValueOf(urlvalue))
				} else {
					mismatchUrlParameters = append(mismatchUrlParameters, fmt.Sprintf("%q (%q -> %q)", metakey, "string", sf.Type.String()))
				}
			} else {
				missingUrlParameters = append(missingUrlParameters, metakey)
			}
		}
	}

	msgs := []string{}
	if len(missingUrlParameters) > 0 {
		msgs = append(msgs, fmt.Sprintf("missing url parameters: %s", strings.Join(missingUrlParameters, ", ")))
	}
	if len(mismatchUrlParameters) > 0 {
		msgs = append(msgs, fmt.Sprintf("type mismatch for url parameters: %s", strings.Join(mismatchUrlParameters, ", ")))
	}
	if len(missingCookies) > 0 {
		msgs = append(msgs, fmt.Sprintf("missing cookies: %s", strings.Join(missingCookies, ", ")))
	}
	if len(mismatchCookies) > 0 {
		msgs = append(msgs, fmt.Sprintf("type mismatch for cookies: %s", strings.Join(mismatchCookies, ", ")))
	}
	if len(msgs) > 0 {
		return fmt.Errorf(strings.Join(msgs, "; "))
	}

	return nil
}

func ParseRequest[Request any](rq *http.Request) (bq *Request, err error) {
	bq = new(Request)
	if err := json.NewDecoder(rq.Body).Decode(bq); err != nil && errors.Is(err, io.ErrUnexpectedEOF) {
		return bq, fmt.Errorf("parsing the request body: %w", err)
	}
	if err := parseRequestByStructFields(rq, bq); err != nil {
		return bq, err
	}
	return
}
