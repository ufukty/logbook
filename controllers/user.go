package controllers

import (
	"fmt"
	"log"
	"logbook/main/crypto"
	"logbook/main/database"
	"logbook/main/parameters"
	"logbook/main/responder"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
)

const (
	PUB_ERR_INVALID_PARAMETERS = "INVALID_PARAMETERS"
	PUB_ERR_INTERNAL_SERVER    = "INTERNAL_SERVER"
	PUB_ERR_INVALID_EMAIL      = "INVALID_EMAIL"
	PUB_ERR_TRY_AGAIN_LATER    = "TRY_AGAIN_LATER"
)

type ErrorDetails struct {
	StatusCode int
	ErrorHint  string
	Error      error
}

var errorVars = map[string]ErrorDetails{
	"UserCreate()/Sanitazation":                 {ErrorHint: PUB_ERR_INVALID_PARAMETERS, StatusCode: http.StatusBadRequest},
	"UserCreate()/SaltCreation":                 {ErrorHint: PUB_ERR_INTERNAL_SERVER, StatusCode: http.StatusInternalServerError},
	"UserCreate()/HashComputation":              {ErrorHint: PUB_ERR_INTERNAL_SERVER, StatusCode: http.StatusInternalServerError},
	"UserCreate()/UserCreation/EmailUniqueness": {ErrorHint: PUB_ERR_INVALID_PARAMETERS, StatusCode: http.StatusBadRequest},
	"UserCreate()/UserCreation/GeneralError":    {ErrorHint: PUB_ERR_INTERNAL_SERVER, StatusCode: http.StatusInternalServerError},
	"UserCreate()/TaskCreateOperation":          {ErrorHint: PUB_ERR_INTERNAL_SERVER, StatusCode: http.StatusInternalServerError},
	"UserCreate()/TaskCreation":                 {ErrorHint: PUB_ERR_INTERNAL_SERVER, StatusCode: http.StatusInternalServerError},
	"UserCreate()/Bookmark Creation":            {ErrorHint: PUB_ERR_INTERNAL_SERVER, StatusCode: http.StatusInternalServerError},
}

func CallErrorHandler(w http.ResponseWriter, r *http.Request, errorObj error, errorPointStamp string) {
	if val, ok := errorVars[errorPointStamp]; ok {
		responder.ErrorHandler(w, r, val.StatusCode, val.ErrorHint, errors.Wrap(errorObj, errorPointStamp))
	} else {
		log.Fatalf("Check the code. 'errorVars' don't have a key for error-point-stamp = '%s'\n", errorPointStamp)
	}
}

/*
 * Objectives for this function
 * DONE: Sanitize user input
 * DONE: Produce unique salt and hash user password with it
 * TODO: Secure against timing-attacks
 * TODO: Check anti-CSRF token
 * DONE: Check account duplication (attempt to register with same e-mail twice)
 * DONE: Create first task
 * DONE: Create privilege table record for created task
 * DONE: Create operation table record for task creation
 * DONE: Create first bookmark
 * TODO: Wrap creation of user-task-bookmark with transaction, rollback on failure to not-lock person to re-register with same email
 */
func UserCreate(w http.ResponseWriter, r *http.Request) {
	params := parameters.UserCreate{}
	if err := params.InputSanitizer(r); err != nil {
		CallErrorHandler(w, r, err, "UserCreate()/Sanitazation")
		return
	}

	// SALT + HASH

	salt, err := crypto.NewSalt()
	if err != nil {
		CallErrorHandler(w, r, err, "UserCreate()/SaltCreation")
		return
	}
	hashedPassword, err := crypto.Argon2Hash([]byte(params.Request.Password), salt)
	if err != nil {
		CallErrorHandler(w, r, err, "UserCreate()/HashComputation")
		return
	}

	// CREATE USER

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

	// CREATE ROOT TASK

	operation := database.Operation{
		UserId:  user.UserId,
		Summary: database.TASK_CREATE,
		Status:  database.SERVER_ORIGINATED,
	}
	result = database.Db.Create(&operation)
	if result.Error != nil {
		CallErrorHandler(w, r, result.Error, "UserCreate()/TaskCreateOperation")
		return
	}
	task := database.Task{
		RevisionId:            operation.OperationId,
		OriginalCreatorUserId: user.UserId,
		ResponsibleUserId:     user.UserId,
		Content:               "ROOT_TASK",
		RootTask:              true,
	}
	result = database.Db.Create(&task)
	if result.Error != nil {
		CallErrorHandler(w, r, result.Error, "UserCreate()/TaskCreation")
		return
	}

	// CREATE ROOT BOOKMARK

	bookmark := database.Bookmark{
		UserId:       user.UserId,
		TaskId:       task.TaskId,
		DisplayName:  "ROOT_BOOKMARK",
		RootBookmark: true,
	}
	result = database.Db.Create(&bookmark)
	if result.Error != nil {
		CallErrorHandler(w, r, result.Error, "UserCreate()/Bookmark Creation")
		return
	}

	responder.SuccessHandler(w, r, params.Response)
}
