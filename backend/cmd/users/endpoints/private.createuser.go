package endpoints

import (
	"fmt"

	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"
)

type CreateUserRequest struct {
	Handle    columns.Username         `json:"handle"`
	Firstname columns.HumanName        `json:"firstname"`
	Lastname  columns.HumanName        `json:"lastname"`
	Birthday  transports.HumanBirthday `json:"birthday"`
	Country   transports.Country       `json:"country"`
	Email     columns.Email            `json:"email"`
	Phone     columns.Phone            `json:"phone"`
	Password  string                   `json:"password"`
}

type CreateUserResponse struct {
	Uid columns.UserId `json:"uid"`
}

// TODO: Authorization
// POST
func (p *Private) CreateUser(w http.ResponseWriter, r *http.Request) {
	bq := CreateUserRequest{}

	if err := bq.Parse(r); err != nil {
		p.l.Println(fmt.Errorf("parsing: %w", err))
		http.Error(w, "parsing", 400)
		return
	}

	uid, err := p.a.CreateUser(r.Context())
	if err != nil {
		p.l.Println(fmt.Errorf("a.CreateUser: %w", err))
		http.Error(w, redact(err), http.StatusInternalServerError)
		return
	}

	bs := CreateUserResponse{
		Uid: uid,
	}

	err = bs.Write(w)
	if err != nil {
		p.l.Println(fmt.Errorf("bq.Write: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
