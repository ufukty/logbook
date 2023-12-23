package document

import (
	c "logbook/controller"
	e "logbook/controller/utilities/errors"
	db "logbook/database"
	"net/http"
)

func createExecutor() (db.Document, *e.Error) {
	var (
		document db.Document
		errs     []error
	)

	// create document table record
	document, errs = db.CreateDocument()
	if errs != nil {
		return db.Document{}, e.New(errs)
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
