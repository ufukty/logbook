package endpoints

import (
	"fmt"

	"logbook/internal/web/serialize"
	"logbook/models/columns"
	"logbook/models/transports"
	"net/http"

	"github.com/ufukty/gohandlers/pkg/types/basics"
)

type CreateUserRequest struct {
	Handle    columns.Username         `json:"handle"`
	Firstname columns.HumanName        `json:"firstname"`
	Lastname  columns.HumanName        `json:"lastname"`
	Birthday  transports.HumanBirthday `json:"birthday"`
	Country   transports.Country       `json:"country"`
	Email     columns.Email            `json:"email"`
	Phone     columns.Phone            `json:"phone"`
	Password  basics.String            `json:"password"`
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

	if issues := bq.Validate(); len(issues) > 0 {
		if err := serialize.ValidationIssues(w, issues); err != nil {
			p.l.Println(fmt.Errorf("serializing validation issues: %w", err))
		}
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
