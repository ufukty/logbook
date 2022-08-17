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
	// log.Println(stdout.String())
}

func SendUserCreateRequest(params *parameters.UserCreate) (*responder.ControllerResponseFields, error) {
	w := httptest.NewRecorder()
	r := PrepareJSONRequest(http.MethodPost, "/user", params.Request)
	UserCreate(w, r)
	apiResponseTemplate := responder.ControllerResponseFields{}
	apiResponseTemplate.Resource = params.Response
	err := DecodeJSONResponse(&apiResponseTemplate, w)
	if err != nil {
		return nil, errors.Wrap(err, "UserCreateTwiceRegistration()")
	}
	return &apiResponseTemplate, nil
}

func TestUserCreateMissingParameters(t *testing.T) {
	ResetTestDatabase()
	database.Init(database.TEST_DB_DSN)
	defer database.CloseConnection()

	params := parameters.UserCreate{}
	params.Request.EmailAddress = "unit-test@golang.example.com"
	response, err := SendUserCreateRequest(&params)
	if err != nil {
		t.Error(err)
	}
	if response.Status == http.StatusOK {
		t.Errorf("Unexpected success for missing parameters")
		t.Errorf("%#v", response)
	}
}

func TestUserCreatePerfectCase(t *testing.T) {
	ResetTestDatabase()
	database.Init(database.TEST_DB_DSN)
	defer database.CloseConnection()

	params := parameters.UserCreate{}
	params.Request.EmailAddress = "unit-test@golang.example.com"
	params.Request.NameSurname = "Quality Assurance"
	params.Request.Password = "lorem-ipsum-dolor-sit-amet"
	response, err := SendUserCreateRequest(&params)
	if err != nil {
		t.Error(err)
	}
	if response.Status != http.StatusOK {
		t.Errorf("Unexpected failure")
		t.Errorf("%#v", response)
	}
}

func TestUserCreateTwiceRegistration(t *testing.T) {
	ResetTestDatabase()
	database.Init(database.TEST_DB_DSN)
	defer database.CloseConnection()

	params := parameters.UserCreate{}
	params.Request.EmailAddress = "unit-test@golang.example.com"
	params.Request.NameSurname = "Quality Assurance"
	params.Request.Password = "lorem-ipsum-dolor-sit-amet"
	response, err := SendUserCreateRequest(&params)
	if err != nil {
		t.Error(err)
	}

	if response.Status != http.StatusOK {
		t.Errorf("Unexpected failure at first registration")
		t.Errorf("%#v", response)
	}

	response, err = SendUserCreateRequest(&params)
	if err != nil {
		t.Error(err)
	}
	if response.Status == http.StatusOK {
		t.Errorf("Unexpected success at second registration")
		t.Errorf("%#v", response)
	}
}

func TestUserCreateInvalidAntiCSRFToken(t *testing.T) {
	ResetTestDatabase()
	database.Init(database.TEST_DB_DSN)
	defer database.CloseConnection()

	t.Errorf("IMPLEMENT")
}
