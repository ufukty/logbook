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

func (param *UserCreate) Sanitize(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&param.Request); err != nil {
		return errors.Wrap(err, "decodeContentTypeJSON")
	}
	if err := param.Request.EmailAddress.TypeCheck(); err != nil {
		return errors.Wrap(err, "sanitize EmailAddress")
	}
	if err := param.Request.Password.TypeCheck(); err != nil {
		return errors.Wrap(err, "sanitize Password")
	}
	if err := param.Request.Salt.TypeCheck(); err != nil {
		return errors.Wrap(err, "sanitize Salt")
	}
	return nil
}

func (param *TaskCreate) Sanitize(r *http.Request) error {
	if err := decodeContentTypeJSON(&param.Request, r); err != nil {
		return errors.Wrap(err, "decodeContentTypeJSON")
	}
	if err := param.Request.AuthorizationToken.TypeCheck(); err != nil {
		return errors.Wrap(err, "sanitize AuthorizationToken")
	}
	if err := param.Request.UserId.TypeCheck(); err != nil {
		return errors.Wrap(err, "sanitize UserId")
	}
	if err := param.Request.Content.TypeCheck(); err != nil {
		return errors.Wrap(err, "sanitize Content")
	}
	if err := param.Request.CurrentRevisionId.TypeCheck(); err != nil {
		return errors.Wrap(err, "sanitize CurrentRevisionId")
	}
	if err := param.Request.SuperTaskId.TypeCheck(); err != nil {
		return errors.Wrap(err, "sanitize SuperTaskId")
	}
	return nil
}
