package endpoints

import (
	"fmt"
	"logbook/cmd/account/app"
	"logbook/cmd/account/database"
	"logbook/cmd/account/service"
	"logbook/config/api"
	"mime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	apicfg, err := api.ReadConfig("../../../api.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("reading api config: %w", err))
	}
	srvcfg, err := service.ReadConfig("../testing.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("reading service config: %w", err))
	}

	err = database.RunMigration(srvcfg)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}

	r := httptest.NewRequest(
		apicfg.Gateways.Public.Services.Account.Endpoints.Create.Method,
		string(apicfg.Gateways.Public.Services.Account.Endpoints.Create.Path),
		// strings.NewReader(`{"firstname": "Tiésto","lastname": "McSingleton","email": "test@test.balaasad.com","password": "123456789"}`),
		strings.NewReader(`{
			"firstname": "Tiésto",
			"lastname": "McSingleton",
			"email": "test@test.balaasad.com",
			"password": "123456789"
		}`),
	)
	r.Header.Add("Content-Type", mime.TypeByExtension(".json"))
	w := httptest.NewRecorder()

	q, err := database.New(srvcfg.Database.Dsn)
	if err != nil {
		t.Fatal(fmt.Errorf("opening database connection: %w", err))
	}
	a := app.New(q)
	ep := New(a)

	ep.CreateUser(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatal("status is not ok")
	}

}
