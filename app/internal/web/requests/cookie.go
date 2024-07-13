package requests

import (
	"fmt"
	"logbook/internal/web/validate"
	"net/http"
)

type Cookie[T comparable] struct {
	Value T
	*http.Cookie
}

var ErrEmpty = fmt.Errorf("empty")

func (c Cookie[T]) Validate() error {
	// nil check (comparable is needed because of this)
	if z := new(T); c.Value == *z {
		return ErrEmpty
	}
	var a any = c.Value
	if v, ok := a.(validate.Validator); ok {
		return v.Validate()
	}
	return nil
}
