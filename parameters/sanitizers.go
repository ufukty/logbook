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
	if err := typeCheckAllFields(map[string]interface{}{
		"name_surname":  &(param.Request.NameSurname),
		"email_address": &(param.Request.EmailAddress),
		"password":      &(param.Request.Password),
	}); err != nil {
		return errors.Wrap(err, "InputSanitizer/Sanitization")
	}
	return nil
}

func (param *TaskCreate) InputSanitizer(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&param.Request); err != nil {
		return errors.Wrap(err, "InputSanitizer/Decoding")
	}
	if err := typeCheckAllFields(map[string]interface{}{
		"authorization_token": &(param.Request.AuthorizationToken),
		"user_id":             &(param.Request.UserId),
		"content":             &(param.Request.Content),
		"super_task_id":       &(param.Request.CurrentRevisionId),
		"current_revision_id": &(param.Request.SuperTaskId),
	}); err != nil {
		return errors.Wrap(err, "InputSanitizer/Sanitization")
	}
	return nil
}

func (param *PlacementArrayHierarchical) InputSanitizer(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&param.Request); err != nil {
		return errors.Wrap(err, "InputSanitizer/Decoding")
	}
	if err := typeCheckAllFields(map[string]interface{}{
		"authorization_token": &(param.Request.AuthorizationToken),
		"user_id":             &(param.Request.UserId),
		"root_task_id":        &(param.Request.RootTaskId),
		"offset":              &(param.Request.Limit),
		"limit":               &(param.Request.Offset),
	}); err != nil {
		return errors.Wrap(err, "InputSanitizer/Sanitization")
	}
	return nil

}
