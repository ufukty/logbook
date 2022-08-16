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
	params.Request.EmailAddress = "unit-test@golang.example.com"
	r := PrepareJSONRequest(http.MethodPost, "/user", params.Request)
	UserCreate(w, r)
	apiResponseTemplate := responder.ControllerResponseFields{}
	apiResponseTemplate.Resource = params.Response
	err := DecodeJSONResponse(&apiResponseTemplate, w)
	if err != nil {
		t.Log(errors.Wrap(err, "UserCreateMissingParameters()"))
	}
	log.Printf("%#v", apiResponseTemplate)
	if apiResponseTemplate.Status != http.StatusBadRequest {
		t.Fail()
	}
	if apiResponseTemplate.ErrorHint != "INVALID_PARAMETERS" {
		t.Fail()
	}
}

func UserCreatePerfectCase(t *testing.T) {
	w := httptest.NewRecorder()
	params := parameters.UserCreate{}
	params.Request.EmailAddress = "unit-test@golang.example.com"
	params.Request.NameSurname = "Quality Assurance"
	params.Request.Password = "lorem-ipsum-dolor-sit-amet"
	r := PrepareJSONRequest(http.MethodPost, "/user", params.Request)
	UserCreate(w, r)
	apiResponseTemplate := responder.ControllerResponseFields{}
	apiResponseTemplate.Resource = params.Response
	err := DecodeJSONResponse(&apiResponseTemplate, w)
	if err != nil {
		t.Log(errors.Wrap(err, "UserCreatePerfectCase()"))
	}
	log.Printf("%#v", apiResponseTemplate)
	if apiResponseTemplate.Status != http.StatusOK {
		t.Fail()
	}
	if apiResponseTemplate.ErrorHint != "" {
		t.Fail()
	}
}

func UserCreateTwiceRegistration(t *testing.T) {
	w := httptest.NewRecorder()
	params := parameters.UserCreate{}
	params.Request.EmailAddress = "unit-test@golang.example.com"
	params.Request.NameSurname = "Quality Assurance"
	params.Request.Password = "lorem-ipsum-dolor-sit-amet"
	r := PrepareJSONRequest(http.MethodPost, "/user", params.Request)
	UserCreate(w, r)
	apiResponseTemplate := responder.ControllerResponseFields{}
	apiResponseTemplate.Resource = params.Response
	err := DecodeJSONResponse(&apiResponseTemplate, w)
	if err != nil {
		t.Log(errors.Wrap(err, "UserCreateTwiceRegistration()"))
	}
	log.Printf("%#v", apiResponseTemplate)
	if apiResponseTemplate.Status != http.StatusBadRequest {
		t.Fail()
	}
	if apiResponseTemplate.ErrorHint != "INVALID_EMAIL" {
		t.Fail()
	}
}

func TestUserCreate(t *testing.T) {
	ResetTestDatabase()
	database.Init(database.TEST_DB_DSN)
	defer database.CloseConnection()

	UserCreateMissingParameters(t)
	// UserCreateInvalidAntiCSRFToken(t)
	UserCreatePerfectCase(t)
	UserCreateTwiceRegistration(t)
}
