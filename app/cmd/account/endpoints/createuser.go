package endpoints

import (
	"errors"
	"fmt"
	"log"
	"logbook/cmd/account/app"
	"logbook/cmd/account/database"
	"logbook/internal/web/reqs"
	"logbook/internal/web/validate"
	"net/http"
)

type CreateUserRequest struct {
	Email     database.Email     `json:"email"`
	Password  string             `json:"password"`
	Firstname database.HumanName `json:"firstname"`
	Lastname  database.HumanName `json:"lastname"`
}

func (bq CreateUserRequest) validate() error {
	return validate.All(map[string]validate.Validator{
		"email":     bq.Email,
		"firstname": bq.Firstname,
		"lastname":  bq.Lastname,
	})
}

/*
 * Objectives for this function
 * DONE: Sanitize user input
 * DONE: Produce unique salt and hash user password with it
 * DONE: Secure against timing-attacks
 * TODO: Check anti-CSRF token
 * DONE: Check account duplication (attempt to register with same e-mail twice)
 * TODO: Create first task
 * TODO: Create privilege table record for created task
 * TODO: Create operation table record for task creation
 * TODO: Create first bookmark
 * TODO: Wrap creation of user-task-bookmark with transaction, rollback on failure to not-lock person to re-register with same email
 */
func (e *Endpoints) CreateUser(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[CreateUserRequest](r)
	if err != nil {
		log.Println(fmt.Errorf("binding: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		log.Println(fmt.Errorf("validation: %w", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = e.app.CreateUser(r.Context(), app.RegistrationParameters{
		Firstname: bq.Firstname,
		Lastname:  bq.Lastname,
		Email:     bq.Email,
		Password:  bq.Password,
	})
	if err != nil {
		log.Println(fmt.Errorf("app.Register: %w", err))
		if errors.Is(err, app.ErrEmailExists) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
