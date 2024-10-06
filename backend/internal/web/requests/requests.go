package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"logbook/internal/utils/lines"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

const bodyBufferLimit = 1 << 20 // 1 MB

func isBodyNeeded[Request any](bq *Request) bool {
	t := reflect.TypeOf(bq).Elem()
	n := t.NumField()
	for i := 0; i < n; i++ {
		ft := t.Field(i)
		if _, ok := ft.Tag.Lookup("json"); ok {
			return true
		}
	}
	return false
}

func parseBody[Request any](w http.ResponseWriter, r *http.Request, dst *Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, bodyBufferLimit)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("request body too large")
	}
	defer r.Body.Close()
	if len(body) == 0 {
		return fmt.Errorf("empty body")
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON")
	}
	return nil
}

type StringAssignable interface {
	Set(string) error
}

func parseUrlFragments[T any](src *http.Request, dst T) error {
	t := reflect.TypeOf(dst).Elem()
	v := reflect.ValueOf(dst).Elem()
	fields := v.NumField()
	vars := mux.Vars(src)
	errs := []string{}

	for i := 0; i < fields; i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		fragmentkey, ok := ft.Tag.Lookup("url")
		if !ok {
			continue
		}
		fragmentvalue, ok := vars[fragmentkey]
		if !ok {
			errs = append(errs, fmt.Sprintf("url doesn't contain url fragment %q for %T.%s", fragmentkey, dst, ft.Name))
			continue
		}

		if fv.Kind() == reflect.Ptr && fv.IsNil() {
			fv.Set(reflect.New(fv.Type().Elem())) // init
		}

		sa, ok := fv.Addr().Interface().(StringAssignable)
		if !ok {
			errs = append(errs, fmt.Sprintf("checking if %T.%s.Value is StringAssignable", dst, ft.Name))
			continue
		}

		if err := sa.Set(fragmentvalue); err != nil {
			errs = append(errs, fmt.Sprintf("%T.%s.Set(): %s", dst, ft.Name, err.Error()))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d errors found:\n%s", len(errs), lines.Join(errs, "+ "))
	}
	return nil
}

func parseCookies[Request any](src *http.Request, dst *Request) error {
	t := reflect.TypeOf(dst).Elem()
	v := reflect.ValueOf(dst).Elem()
	fields := v.NumField()

	errs := []string{}
	for i := 0; i < fields; i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		cookiename, ok := ft.Tag.Lookup("cookie")
		if !ok {
			continue
		}
		cookievalue, err := src.Cookie(cookiename)
		if err != nil {
			errs = append(errs, fmt.Sprintf("checking cookies for %q for %T.%s: %s", cookiename, dst, ft.Name, err.Error()))
			continue
		}

		if fv.Kind() == reflect.Ptr && fv.IsNil() {
			fv.Set(reflect.New(fv.Type().Elem())) // init
		}

		ct := ft.Type
		if ct.Kind() == reflect.Struct && strings.HasPrefix(ct.Name(), "Cookie") { // Cookie type is used
			_, okV := ct.FieldByName("Value")  // StringAssignable
			_, okD := ct.FieldByName("Cookie") // http.Cookie
			if !(okV && okD) {
				errs = append(errs, fmt.Sprintf("'Cookie' named struct type doesn't contain either or both 'Value' and '*http.Cookie', used for %T.%s", dst, ft.Name))
				continue
			}

			valueField := fv.FieldByName("Value")
			detailField := fv.FieldByName("Cookie")

			sa, ok := valueField.Addr().Interface().(StringAssignable)
			if !ok {
				errs = append(errs, fmt.Sprintf("checking if %T.%s.Value is StringAssignable", dst, ft.Name))
				continue
			}
			if err := sa.Set(cookievalue.Value); err != nil {
				errs = append(errs, fmt.Sprintf("%T.%s.Value.Set(): %s", dst, ft.Name, err.Error()))
				continue
			}

			httpcookie, ok := detailField.Addr().Interface().(**http.Cookie)
			if !ok {
				errs = append(errs, fmt.Sprintf("checking if the type of %T.%s is *http.Cookie", dst, ft.Name))
				continue
			}
			*httpcookie = cookievalue
		} else {
			errs = append(errs, fmt.Sprintf("unsupported type %s for %T.%s", ct.Kind(), dst, ft.Name))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d errors found:\n%s", len(errs), lines.Join(errs, "+ "))
	}
	return nil
}

func ParseRequest[Request any](w http.ResponseWriter, r *http.Request, dst *Request) error {
	if isBodyNeeded(dst) {
		if err := parseBody(w, r, dst); err != nil {
			return fmt.Errorf("parsing the request body: %w", err)
		}
	}
	if err := parseUrlFragments(r, dst); err != nil {
		return fmt.Errorf("parsing url fragments: %w", err)
	}
	if err := parseCookies(r, dst); err != nil {
		return fmt.Errorf("parsing cookies: %w", err)
	}
	return nil
}
