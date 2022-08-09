package parameters

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func PrepareRequestForContentTypeJSON(method string, target string, json []byte) *http.Request {
	body := bytes.NewBuffer(json)
	r := httptest.NewRequest(method, target, body)
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	return r
}

func TestUserCreate(t *testing.T) {
	t.Log("should PASS")
	r := PrepareRequestForContentTypeJSON(http.MethodPost, "/user", []byte(`{
		"email_address": "testUserCreate@golang.example.com",
		"salt": "loremipsumdolorsitamet",
		"password": "123456789"
	}`))
	parameters := UserCreate{}
	if err := parameters.Sanitize(r); err != nil {
		t.Errorf(err.Error())
	} else {
		t.Log(err)
		t.Log(parameters)
	}
	t.Log("should FAIL")
	r = PrepareRequestForContentTypeJSON(http.MethodPost, "/user", []byte(`{
		"email_address": "testUserCreate@golang.example.com"
	}`))
	parameters = UserCreate{}
	if err := parameters.Sanitize(r); err == nil {
		t.Errorf("Misconfigured input has accepted.")
	} else {
		t.Log(err)
		t.Log(parameters)
	}
}

func TestTaskCreate(t *testing.T) {
	t.Log("should PASS")
	r := PrepareRequestForContentTypeJSON(http.MethodPost, "/task/create", []byte(`{
		"authorization_token": "TODO: implement auth later",
		"user_id": "00000000-0000-0000-0000-000000000000",
		"super_task_id": "00000000-0000-0000-0000-000000000000",
		"current_revision_id": "00000000-0000-0000-0000-000000000000",
		"content": "First task"
	}`))
	parameters := TaskCreate{}
	if err := parameters.Sanitize(r); err != nil {
		t.Errorf(err.Error())
	} else {
		t.Log(err)
		t.Log(parameters)
	}
	t.Log("should FAIL")
	r = PrepareRequestForContentTypeJSON(http.MethodPost, "/task/create", []byte(`{
		"content": "First task"
	}`))
	parameters = TaskCreate{}
	if err := parameters.Sanitize(r); err == nil {
		t.Errorf(err.Error())
	} else {
		t.Log(err)
		t.Log(parameters)
	}
}
