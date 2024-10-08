package endpoints

import (
	"context"
	"fmt"
	"logbook/cmd/account/api/public/app"
	"logbook/cmd/account/database"
	"logbook/cmd/account/service"
	"logbook/config/api"
	"logbook/internal/logger"
	"logbook/models"
	"mime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MockInstanceSource []models.Instance

func (m *MockInstanceSource) Instances() ([]models.Instance, error) {
	return *m, nil
}

func TestCreateUser(t *testing.T) {
	apicfg, err := api.ReadConfig("../../../api.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("reading api config: %w", err))
	}
	srvcfg, err := service.ReadConfig("../local.yml")
	if err != nil {
		t.Fatal(fmt.Errorf("reading service config: %w", err))
	}

	err = database.RunMigration(srvcfg)
	if err != nil {
		t.Fatal(fmt.Errorf("running migration: %w", err))
	}

	r := httptest.NewRequest(
		apicfg.Public.Services.Account.Endpoints.CreateUser.Method,
		apicfg.Public.Services.Account.Endpoints.CreateUser.Path,
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

	pool, err := pgxpool.New(context.Background(), srvcfg.Database.Dsn)
	if err != nil {
		t.Fatal(fmt.Errorf("opening database connection: %w", err))
	}
	defer pool.Close()

	a := app.New(pool, apicfg, nil) // FIXME: mock objectives service?
	ep := New(a, logger.New("test"))

	ep.CreateUser(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatal("status is not ok")
	}

}
