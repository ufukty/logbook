package parameters

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/pkg/errors"
)

func PrepareRequestForContentTypeJSON(method string, target string, json []byte) *http.Request {
	body := bytes.NewBuffer(json)
	r := httptest.NewRequest(method, target, body)
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	return r
}

func DecodeContentTypeJSON(param interface{}, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		return errors.Wrap(err, "decodeContentTypeJSON")
	}
	return nil
}

func sanitizeFields(fields []interface{}) error {
	for _, field := range fields {
		if typeCheckableField, ok := field.(TypeCheckable); ok {
			if err := typeCheckableField.TypeCheck(); err != nil {
				return errors.Wrapf(err, "sanitizeFields parameter = <%s>", field)
			}
		}
	}
	return nil
}
