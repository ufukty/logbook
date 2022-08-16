package parameters

import (
	"log"
	"net/http"
	"testing"

	"github.com/pkg/errors"
)

// check for if InputSanitizer works/fails when it is appropriate (for different inputs)
func TestUserCreate(t *testing.T) {
	log.Println("should PASS")
	r := PrepareJSONRequest(http.MethodPost, "/user", []byte(`{
		"email_address": "testUserCreate@golang.example.com",
		"name_surname": "Quality Assurance",
		"password": "123456789"
	}`))
	params := UserCreate{}
	if err := params.InputSanitizer(r); err != nil {
		t.Error(errors.Wrap(err, "TestUserCreate()"))
	} else {
		log.Println(err)
		log.Println(params)
	}
	log.Println("should FAIL")
	r = PrepareJSONRequest(http.MethodPost, "/user", []byte(`{
		"email_address": "testUserCreate@golang.example.com"
	}`))
	params = UserCreate{}
	if err := params.InputSanitizer(r); err == nil {
		t.Errorf("Misconfigured input has accepted.")
	} else {
		log.Println(err)
		log.Println(params)
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
