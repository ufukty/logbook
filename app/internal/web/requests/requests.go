package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"logbook/internal/utilities/slicew/lines"
	"net/http"
	"reflect"
)

// func parseUrlFragments[T any](src *http.Request, dst T) error {
// 	t := reflect.TypeOf(dst).Elem()
// 	v := reflect.ValueOf(dst).Elem()
// 	fields := v.NumField()
// 	vars := mux.Vars(src)
// 	errs := []string{}

// 	for i := 0; i < fields; i++ {
// 		field := t.Field(i)
// 		v := v.Field(i)

// 		fragmentkey, ok := field.Tag.Lookup("url")
// 		if !ok {
// 			continue
// 		}

// 		fragmentvalue, ok := vars[fragmentkey]
// 		if !ok {
// 			errs = append(errs, fmt.Sprintf("url doesn't contain url fragment %q for %T.%s", fragmentkey, dst, field.Name))
// 			continue
// 		}

// 		if v.Kind() == reflect.Ptr && v.IsNil() {
// 			v.Set(reflect.New(v.Type().Elem())) // init
// 		}

// 		c, ok := v.Addr().Interface().(urlholder)
// 		if !ok {
// 			errs = append(errs, fmt.Sprintf("asserting %T.%s is value assignable from string", dst, field.Name))
// 			continue
// 		}

// 		if err := c.cookievalue(fragmentvalue); err != nil {
// 			errs = append(errs, fmt.Sprintf("assigning url fragment value (%s: %s) to %T.%s: %s", fragmentkey, fragmentvalue, dst, field.Name, err.Error()))
// 		}
// 	}

// 	if len(errs) > 0 {
// 		return fmt.Errorf("%d errors found:\n%s", len(errs), lines.Join(errs, "+ "))
// 	}
// 	return nil
// }

func parseCookies[T any](src *http.Request, dst T) error {
	t := reflect.TypeOf(dst).Elem()
	v := reflect.ValueOf(dst).Elem()
	fields := v.NumField()

	errs := []string{}
	for i := 0; i < fields; i++ {
		v := v.Field(i)
		field := t.Field(i)

		cookiename, ok := field.Tag.Lookup("cookie")
		if !ok {
			continue
		}

		cookievalue, err := src.Cookie(cookiename)
		if err != nil {
			errs = append(errs, fmt.Sprintf("checking cookies for %q for %T.%s: %s", cookiename, dst, field.Name, err.Error()))
			continue
		}

		if v.Kind() == reflect.Ptr && v.IsNil() {
			v.Set(reflect.New(v.Type().Elem())) // init
		}

		ch, ok := v.Addr().Interface().(cookieholder)
		if !ok {
			errs = append(errs, fmt.Sprintf("asserting %T.%s is value assignable from string", dst, field.Name))
			continue
		}

		if err := ch.setCookieValue(cookievalue.Value); err != nil {
			errs = append(errs, fmt.Sprintf("assigning cookie value to %T.%s: %s", dst, field.Name, err.Error()))
			continue
		}
		ch.setCookieDetails(cookievalue)
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d errors found:\n%s", len(errs), lines.Join(errs, "+ "))
	}
	return nil
}

func ParseRequest[Request any](rq *http.Request) (bq *Request, err error) {
	bq = new(Request)
	if err := json.NewDecoder(rq.Body).Decode(bq); err != nil && errors.Is(err, io.ErrUnexpectedEOF) {
		return bq, fmt.Errorf("parsing the request body: %w", err)
	}
	// if err := parseUrlFragments(rq, bq); err != nil {
	// 	return bq, fmt.Errorf("parsing url fragments: %w", err)
	// }
	if err := parseCookies(rq, bq); err != nil {
		return bq, fmt.Errorf("parsing cookies: %w", err)
	}
	return
}
