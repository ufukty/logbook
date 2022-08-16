package controllers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"logbook/main/database"
	"logbook/main/parameters"
	responder "logbook/main/responder"

	"github.com/pkg/errors"
)

func LoadTestDatabase() {
	database.Init(database.TEST_DB_DSN)
}

func ResetTestDatabase() {
	cmd := exec.Command("make", "migrate")
	cmd.Dir = "../"

	// cmd.Stdin = strings.NewReader("and old falcon")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	err := cmd.Run()

	if err != nil {
		log.Fatalln(stdout.String(), err)
	}
	log.Println(stdout.String())
}

func UserCreateMissingParameters(t *testing.T) {
	w := httptest.NewRecorder()
	params := parameters.UserCreate{}
	params.Request.EmailAddress = "testUserCreate@golang.example.com"
	r := parameters.PrepareJSONRequest(http.MethodPost, "/user", params.Request)
	UserCreate(w, r)
	apiResponseTemplate := responder.ControllerResponseFields{}
	apiResponseTemplate.Resource = params.Response
	err := parameters.DecodeJSONResponse(&apiResponseTemplate, w)
	if err != nil {
		log.Println(errors.Wrap(err, "UserCreateMissingParameters()"))
		t.Fail()
	}
	log.Printf("%#v", apiResponseTemplate)
	if apiResponseTemplate.Status != http.StatusBadRequest {
		t.Fail()
	}
}

// test happy path
// test failing conditions:
// - missing parameters
// - invalid ant-CSRF token
// - second registration -> same username / e-mail
func TestUserCreate(t *testing.T) {
	ResetTestDatabase()
	database.Init(database.TEST_DB_DSN)
	defer database.CloseConnection()

	UserCreateMissingParameters(t)
	// UserCreateInvalidAntiCSRFToken(t)
	// UserCreateTwiceRegistration(t)
	// UserCreatePerfectCase(t)
	// SHOULD FAIL: because empty/invalid ant-CSRF token

	// SHOULD PASS:
	// r = parameters.PrepareJSONRequest(http.MethodPost, "/user", []byte(`{
	// 	"email_address": "testUserCreate@golang.example.com",
	// 	"random_number": "loremipsumdolorsitamet",
	// 	"password": "123456789"
	// }`))
	// UserCreate(w, r)

	// SHOULD FAIL: Register twice

	// parameters.DecodeContentTypeJSON(parameters
}
