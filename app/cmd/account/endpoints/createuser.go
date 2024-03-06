package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/account/database"
	"logbook/internal/web/reqs"
	"logbook/internal/web/validate"
	"net/http"

	"github.com/alexedwards/argon2id"
)

type CreateUserRequest struct {
	Email       database.Email     `json:"email"`
	NameSurname database.HumanName `json:"name_surname"`
	Password    string             `json:"password"`
	Username    database.Username  `json:"username"`
}

func (bq CreateUserRequest) validate() error {
	return validate.All(map[string]validate.Validator{
		"email":        bq.Email,
		"name surname": bq.NameSurname,
		"username":     bq.Username,
	})
}

type CreateUserResponse struct{}

var argon2idParams = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

/*
 * Objectives for this function
 * TODO: Sanitize user input
 * TODO: Produce unique salt and hash user password with it
 * TODO: Secure against timing-attacks
 * TODO: Check anti-CSRF token
 * TODO: Check account duplication (attempt to register with same e-mail twice)
 * TODO: Create first task
 * TODO: Create privilege table record for created task
 * TODO: Create operation table record for task creation
 * TODO: Create first bookmark
 * TODO: Wrap creation of user-task-bookmark with transaction, rollback on failure to not-lock person to re-register with same email
 */
func (e *Endpoints) CreateUser(w http.ResponseWriter, r *http.Request) {
	bq, err := reqs.ParseRequest[CreateUserRequest](r)
	if err != nil {
		log.Println(fmt.Errorf("parsing request: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := bq.validate(); err != nil {
		log.Println(fmt.Errorf("validating request parameters: %w", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	hash, err := argon2id.CreateHash(bq.Password, argon2idParams)
	if err != nil {
		log.Println(fmt.Errorf("creating hash: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}
