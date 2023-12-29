package endpoints

import (
	"fmt"
	"log"
	"logbook/cmd/user/database"
	"logbook/internal/web/crypto"
	"logbook/internal/web/reqs"
	"net/http"
)

type CreateUserRequest struct {
	NameSurname NonEmptyString    `json:"name_surname"`
	Username    database.Username `json:"username"`
	Password    string            `json:"password"`
}

func (bq CreateUserRequest) validate() error {
	if !bq.Username.Validate() {
		return fmt.Errorf("invalid value for username parameter")
	}
	return nil
}

type CreateUserResponse struct{}

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

	salt, err := crypto.NewSalt()
	if err != nil {
		log.Println(fmt.Errorf("creating salt: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	hash, err := crypto.Argon2Hash([]byte(bq.Password), salt)
	if err != nil {
		log.Println(fmt.Errorf("creating hash: %w", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	user := database.User{
		NameSurname:       string(params.Request.NameSurname),
		EmailAddress:      string(params.Request.EmailAddress),
		SaltBase64Encoded: string(crypto.Base64Encode(salt)),
		HashEncoded:       hashedPassword,
	}
	result := database.Db.Create(&user)
	if result.Error != nil {
		if errorCode := database.StripSQLState(fmt.Sprint(result.Error)); errorCode == pgerrcode.UniqueViolation {
			CallErrorHandler(w, r, err, "UserCreate()/UserCreation/EmailUniqueness")
		} else {
			CallErrorHandler(w, r, err, "UserCreate()/UserCreation/GeneralError")
		}
		return
	}
}
