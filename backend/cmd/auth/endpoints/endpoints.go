package endpoints

import (
	"logbook/cmd/auth/app"
	"net/http"
)

type Endpoints struct {
	App *app.App
}

func (eps *Endpoints) LoginPage(w *http.ResponseWriter, r *http.Request) {

}

func (eps *Endpoints) SignupPage(w *http.ResponseWriter, r *http.Request) {

}
