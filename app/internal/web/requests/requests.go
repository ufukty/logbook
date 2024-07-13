package requests

import (
	"encoding/json"
	"fmt"
	"logbook/internal/utilities/slicew/lines"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

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

func ParseRequest[Request any](src *http.Request, dst *Request) error {
	if err := json.NewDecoder(src.Body).Decode(dst); err != nil {
		return fmt.Errorf("parsing the request body: %w", err)
	}
	if err := parseUrlFragments(src, dst); err != nil {
		return fmt.Errorf("parsing url fragments: %w", err)
	}
	if err := parseCookies(src, dst); err != nil {
		return fmt.Errorf("parsing cookies: %w", err)
	}
	return nil
}
