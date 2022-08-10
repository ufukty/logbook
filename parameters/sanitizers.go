package parameters

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func (param *UserCreate) InputSanitizer(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&param.Request); err != nil {
		return errors.Wrap(err, "InputSanitizer/Decoding")
	}
	if err := sanitizeFields([]interface{}{
		&(param.Request.EmailAddress),
		&(param.Request.Password),
		&(param.Request.Salt),
	}); err != nil {
		return errors.Wrap(err, "InputSanitizer/Sanitization")
	}
	return nil
}

func (param *TaskCreate) InputSanitizer(r *http.Request) error {
	if err := decodeContentTypeJSON(&param.Request, r); err != nil {
		return errors.Wrap(err, "InputSanitizer/Decoding")
	}
	if err := sanitizeFields([]interface{}{
		&(param.Request.AuthorizationToken),
		&(param.Request.UserId),
		&(param.Request.Content),
		&(param.Request.CurrentRevisionId),
		&(param.Request.SuperTaskId),
	}); err != nil {
		return errors.Wrap(err, "InputSanitizer/Sanitization")
	}
	return nil
}

func (param *PlacementArrayHierarchical) InputSanitizer(r *http.Request) error {
	if err := decodeContentTypeJSON(&param.Request, r); err != nil {
		return errors.Wrap(err, "InputSanitizer/Decoding")
	}
	if err := sanitizeFields([]interface{}{
		&(param.Request.AuthorizationToken),
		&(param.Request.UserId),
		&(param.Request.RootTaskId),
		&(param.Request.Limit),
		&(param.Request.Offset),
	}); err != nil {
		return errors.Wrap(err, "InputSanitizer/Sanitization")
	}
	return nil

}