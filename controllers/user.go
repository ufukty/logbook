package controllers

import (
	"fmt"
	database "logbook/main/database"
	"logbook/main/parameters"
	"logbook/main/responder"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
	"github.com/tvdburgt/go-argon2"
)

var defaultArgon2Config = &argon2.Context{
	Iterations:  3,
	Memory:      1 << 12, // 4 MiB
	Parallelism: 1,
	HashLen:     32,
	Mode:        argon2.ModeArgon2id,
	Version:     argon2.Version13,
}

func VerifyHash(storedHash, clearText string) (bool, error) {
	result, err := argon2.VerifyEncoded(storedHash, []byte(clearText))
	if err != nil {
		return false, errors.Wrap(err, "VerifyHash()")
	}
	return result, nil
}

func Hash(clear, salt string) (string, error) {
	hashedString, err := argon2.HashEncoded(defaultArgon2Config, []byte(clear), []byte(salt))
	if err != nil {
		return "", errors.Wrap(err, "Hash()")
	}
	return hashedString, nil
}

// Objectives for this function
// * Sanitize user input
// TODO: * Secure against timing-attacks
// TODO: * Check anti-CSRF token
// TODO: * Check account duplication (attempt to register with same e-mail twice)
// TODO: * Create first task
// TODO: * Create first bookmark
func UserCreate(w http.ResponseWriter, r *http.Request) {
	params := parameters.UserCreate{}
	if err := params.InputSanitizer(r); err != nil {
		responder.ErrorHandler(w, r, http.StatusBadRequest, "INVALID_PARAMETERS", errors.Wrap(err, "UserCreate"))
		return
	}

	hashedPassword, err := Hash(string(params.Request.Password), string(params.Request.RandomNumber))
	if err != nil {
		responder.ErrorHandler(w, r, http.StatusInternalServerError, "INVALID_PARAMETERS", errors.Wrap(err, "UserCreate()"))
		return
	}

	user := database.User{
		NameSurname:    string(params.Request.NameSurname),
		EmailAddress:   string(params.Request.EmailAddress),
		Salt:           string(params.Request.RandomNumber),
		HashedPassword: hashedPassword,
	}

	// create db record for new user
	result := database.Db.Create(&user)
	if result.Error != nil {
		errorCode := database.StripSQLState(fmt.Sprint(result.Error))
		if errorCode == pgerrcode.UniqueViolation {
			responder.ErrorHandler(w, r, http.StatusBadRequest, "INVALID_EMAIL", errors.Wrap(err, "UserCreate()"))
			return
		}
		responder.ErrorHandler(w, r, http.StatusInternalServerError, "INVALID_PARAMETERS", errors.Wrap(err, "UserCreate()"))
		return
	}

	responder.SuccessHandler(w, r, params.Response)
}
