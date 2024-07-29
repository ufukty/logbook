package endpoints

import (
	"fmt"
	"logbook/cmd/account/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/validate"
	database "logbook/models/columns"
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
	bq := &CreateUserRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		e.l.Println(fmt.Errorf("binding: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		e.l.Println(fmt.Errorf("validation: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	err := e.app.CreateUser(r.Context(), app.RegistrationParameters{
		Firstname: bq.Firstname,
		Lastname:  bq.Lastname,
		Email:     bq.Email,
		Password:  bq.Password,
	})
	if err != nil {
		e.l.Println(fmt.Errorf("app.Register: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
