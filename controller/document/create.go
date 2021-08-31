package document

import (
	c "logbook/main/controller"
	db "logbook/main/database"
	"net/http"
)

func createExecutor() (db.Document, []error) {
	var (
		document db.Document
		errs     []error
	)

	// create document table record
	document, errs = db.CreateDocumentWithTaskGroups(db.Document{})
	if errs != nil {
		return db.Document{}, append(errs, c.ErrCreateDocumentCreateCreationPhase)
	}

	taskGroups, errs := db.GetTaskGroupsByDocumentId(document.DocumentId)
	if errs != nil {
		return db.Document{}, append(errs, c.ErrCreateDocumentSelectionPhase)
	}

	document.TaskGroups = taskGroups
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
