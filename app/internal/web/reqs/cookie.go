package reqs

import (
	"fmt"
	"logbook/internal/web/validate"
	"net/http"
	"reflect"
	"strconv"
)

type Cookie[T comparable] struct {
	Value T
	*http.Cookie
}

var ErrEmpty = fmt.Errorf("empty")

func (c Cookie[T]) Validate() error {
	if z := new(T); c.Value == *z { // so, below won't panic
		return ErrEmpty
	}
	var a any = c.Value
	if v, ok := a.(validate.Validator); ok {
		return v.Validate()
	}
	return nil
}

func (c *Cookie[T]) setCookieValue(value string) error {
	v := reflect.ValueOf(c).Elem()
	fv := v.FieldByName("Value")

	// rules are based on https://github.com/syntaqx/cookie/blob/c9e46f45600911422b1c3742be5ae5f9f87b70d3/populate.go
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, fv.Type().Bits())
		if err != nil {
			panic(err)
		}
		fv.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(value, 10, fv.Type().Bits())
		if err != nil {
			return err
		}
		fv.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, fv.Type().Bits())
		if err != nil {
			return err
		}
		fv.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		fv.SetBool(b)
	default:
		return fmt.Errorf("unsupported type for source (%q) or destination (%T)", fv, fv.Kind())
	}

	return nil
}

func (c *Cookie[T]) setCookieDetails(cookie *http.Cookie) {
	c.Cookie = cookie
}
