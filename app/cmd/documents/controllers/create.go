package document

import (
	"logbook/cmd/tasks/database"
	c "logbook/controller"
	e "logbook/controller/utilities/errors"
	"net/http"
)

func createExecutor() (database.Document, *e.Error) {
	var (
		document database.Document
		errs     []error
	)

	// create document table record
	document, errs = database.CreateDocument()
	if errs != nil {
		return database.Document{}, e.New(errs)
	}

	return document, nil
}

func Create(w http.ResponseWriter, r *http.Request) {
	// ipAddress := (*r).RemoteAddr
	// userAgent := (*r).Header.Get("User-Agent")
	document, errs := createExecutor()
	if errs != nil {
		c.ErrorHandler(w, r, errs)
	} else {
		c.SuccessHandler(w, r, document)
	}
}
