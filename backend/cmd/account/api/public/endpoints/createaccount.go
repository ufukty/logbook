package endpoints

import (
	"fmt"
	"logbook/cmd/account/api/public/app"
	"logbook/internal/web/requests"
	"logbook/internal/web/router/reception"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"net/http"
)

type CreateAccountRequest struct {
	Email     columns.Email     `json:"email"`
	Password  string            `json:"password"`
	Firstname columns.HumanName `json:"firstname"`
	Lastname  columns.HumanName `json:"lastname"`
}

func (bq CreateAccountRequest) validate() error {
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
func (e *Endpoints) CreateAccount(id reception.RequestId, store *reception.Store, w http.ResponseWriter, r *http.Request) error {
	bq := &CreateAccountRequest{}

	if err := requests.ParseRequest(w, r, bq); err != nil {
		http.Error(w, redact(err), http.StatusBadRequest)
		return fmt.Errorf("binding: %w", err)
	}

	if err := bq.validate(); err != nil {
		http.Error(w, redact(err), http.StatusBadRequest)
		return fmt.Errorf("validation: %w", err)
	}

	err := e.a.CreateAccount(r.Context(), app.CreateAccountRequest{
		Firstname: bq.Firstname,
		Lastname:  bq.Lastname,
		Email:     bq.Email,
		Password:  bq.Password,
	})
	if err != nil {
		http.Error(w, redact(err), http.StatusInternalServerError)
		return fmt.Errorf("app.Register: %w", err)
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
