package parameters

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/pkg/errors"
)

func PrepareJSONRequest(method string, target string, json []byte) *http.Request {
	body := bytes.NewBuffer(json)
	r := httptest.NewRequest(method, target, body)
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	return r
}

func typeCheckAllFields(fields map[string]interface{}) error {
	for handle, field := range fields {
		if typeCheckableField, ok := field.(TypeCheckable); ok {
			if err := typeCheckableField.TypeCheck(); err != nil {
				return errors.Wrapf(err, "typeCheckAllFields %s", handle)
			}
		}
	}
	return nil
}
