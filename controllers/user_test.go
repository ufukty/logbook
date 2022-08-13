package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"logbook/main/database"
	"logbook/main/parameters"
)

func TestUserCreate(t *testing.T) {
	// test happy path

	// test failing conditions
	//   missing parameters
	//     dsd
	//// second registration -> same username / e-mail

	database.Init(database.TEST_DB_DSN)

	randomEmail := fmt.Sprintf("test_user_%s", rand.Int())
	randomNumber := rand.Int()

	r := parameters.PrepareRequestForContentTypeJSON(http.MethodPost, "/user", []byte(`{
		"email_address": "testUserCreate@golang.example.com",
		"random_number": "loremipsumdolorsitamet",
		"password": "123456789"
	}`))
	w := httptest.NewRecorder()
	UserCreate(w, r)
	// res := w.Result()
	// defer res.Body.Close()
	// data, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	t.Errorf("expected error to be nil got %v", err)
	// }
	// if string(data) != "ABC" {
	// 	t.Errorf("expected ABC got %v", string(data))
	// }
}
