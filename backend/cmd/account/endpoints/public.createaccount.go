package endpoints

import (
	"fmt"
	"logbook/cmd/account/app"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
)

type CreateAccountRequest struct {
	Email     columns.Email       `json:"email"`
	Password  transports.Password `json:"password"`
	Firstname columns.HumanName   `json:"firstname"`
	Lastname  columns.HumanName   `json:"lastname"`
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
// POST
func (p *Public) CreateAccount(w http.ResponseWriter, r *http.Request) {
	bq := &CreateAccountRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("binding: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	if err := validate.RequestFields(bq); err != nil {
		p.l.Println(fmt.Errorf("validation: %w", err))
		http.Error(w, redact(err), http.StatusBadRequest)
		return
	}

	err := p.a.CreateAccount(r.Context(), app.CreateAccountRequest{
		Firstname: bq.Firstname,
		Lastname:  bq.Lastname,
		Email:     bq.Email,
		Password:  bq.Password,
	})
	if err != nil {
		p.l.Println(fmt.Errorf("app.Register: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
