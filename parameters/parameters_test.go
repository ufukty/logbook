package parameters

import (
	"net/http"
	"testing"
)

func TestUserCreate(t *testing.T) {
	t.Log("should PASS")
	r := PrepareJSONRequest(http.MethodPost, "/user", []byte(`{
		"email_address": "testUserCreate@golang.example.com",
		"salt": "loremipsumdolorsitamet",
		"password": "123456789"
	}`))
	parameters := UserCreate{}
	if err := parameters.InputSanitizer(r); err != nil {
		t.Errorf(err.Error())
	} else {
		t.Log(err)
		t.Log(parameters)
	}
	t.Log("should FAIL")
	r = PrepareJSONRequest(http.MethodPost, "/user", []byte(`{
		"email_address": "testUserCreate@golang.example.com"
	}`))
	parameters = UserCreate{}
	if err := parameters.InputSanitizer(r); err == nil {
		t.Errorf("Misconfigured input has accepted.")
	} else {
		t.Log(err)
		t.Log(parameters)
	}
}

func TestTaskCreate(t *testing.T) {
	t.Log("should PASS")
	r := PrepareJSONRequest(http.MethodPost, "/task/create", []byte(`{
		"authorization_token": "TODO: implement auth later",
		"user_id": "00000000-0000-0000-0000-000000000000",
		"super_task_id": "00000000-0000-0000-0000-000000000000",
		"current_revision_id": "00000000-0000-0000-0000-000000000000",
		"content": "First task"
	}`))
	parameters := TaskCreate{}
	if err := parameters.InputSanitizer(r); err != nil {
		t.Errorf(err.Error())
	} else {
		t.Log(err)
		t.Log(parameters)
	}
	t.Log("should FAIL")
	r = PrepareJSONRequest(http.MethodPost, "/task/create", []byte(`{
		"content": "First task"
	}`))
	parameters = TaskCreate{}
	if err := parameters.InputSanitizer(r); err == nil {
		t.Errorf(err.Error())
	} else {
		t.Log(err)
		t.Log(parameters)
	}
}
