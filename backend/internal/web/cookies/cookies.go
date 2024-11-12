package cookies

import (
	"fmt"
	"net/http"
)

type c interface {
	Validate() error
	FromCookie(string) error
	ToCookie() (string, error)
}

type Must[C c] struct {
	*http.Cookie   // to preserve access to cookie attributes...
	Value        C // ...while providing seamless access to type checked cookie value
}

func (c *Must[T]) FromCookie(h *http.Cookie) error {
	if h == nil {
		return fmt.Errorf("cookie not found")
	}
	c.Cookie = h
	if h.Value == "" {
		return fmt.Errorf("empty value")
	}
	return c.Value.FromCookie(h.Value)
}

func (c *Must[T]) SetCookie(w http.ResponseWriter) {
	http.SetCookie(w, c.Cookie)
}

func (c Must[T]) Validate() error {
	return c.Value.Validate()
}

type Optional[C c] struct {
	*http.Cookie
	Value C
}

func (c *Optional[T]) FromCookie(h *http.Cookie) error {
	if h == nil {
		return nil
	}
	if h.Value == "" {
		return nil
	}
	c.Cookie = h
	return c.Value.FromCookie(h.Value)
}

func (c *Optional[T]) SetCookie(w http.ResponseWriter) {
	http.SetCookie(w, c.Cookie)
}

func (c Optional[T]) Validate() error {
	return c.Value.Validate()
}
