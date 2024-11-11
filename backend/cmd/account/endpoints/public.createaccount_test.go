package endpoints

import (
	"context"
	"fmt"
	"logbook/cmd/account/app"
	"logbook/cmd/account/database"
	"logbook/cmd/account/service"
	objectives "logbook/cmd/objectives/client"
	"logbook/internal/startup"
	"mime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCreateAccount(t *testing.T) {
	l, srvcfg, _, err := startup.TestDependenciesWithServiceConfig("account", service.ReadConfig)
	if err != nil {
		t.Fatal(fmt.Errorf("startup.TestDependenciesWithServiceConfig: %w", err))
	}

	err = database.RunMigration(srvcfg)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}

	r := httptest.NewRequest("DUMMY", "/DUMMY", strings.NewReader(`{
		"firstname": "Tiésto",
		"lastname": "McSingleton",
		"email": "test@test.balaasad.com",
		"password": "123456789"
	}`))
	r.Header.Add("Content-Type", mime.TypeByExtension(".json"))
	w := httptest.NewRecorder()

	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		t.Fatal(fmt.Errorf("opening database connection: %w", err))
	}
	defer pool.Close()

	a := app.New(pool, &objectives.Mock{})
	pu := NewPublic(a, l)

	pu.CreateAccount(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatal("status is not ok")
	}

}
