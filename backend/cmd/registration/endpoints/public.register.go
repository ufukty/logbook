package endpoints

import (
	"fmt"
	"logbook/cmd/registration/app"
	"logbook/internal/web/validate"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
	"time"
)

type CreateAccountRequest struct {
	CsrfToken string

	Firstname columns.HumanName `json:"firstname"`
	Lastname  columns.HumanName `json:"lastname"`
	Birthday  time.Time
	// Country   columns.Country

	Email columns.Email `json:"email"`
	// EmailGrant columns.EmailGrant

	// Phone      columns.Phone
	// PhoneGrant columns.PhoneGrant

	Password transports.Password `json:"password"`
}

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

	err := p.a.Register(r.Context(), app.RegisterRequest{
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
