package parameters

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func decodeContentTypeJSON(param interface{}, r *http.Request) error {
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
				return errors.Wrap(err, "sanitizeFields")
			}
		}
	}
	return nil
}
