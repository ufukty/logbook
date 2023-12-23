package task

import (
	c "logbook/controller"
	e "logbook/controller/utilities/errors"
	v "logbook/controller/utilities/validate"
	db "logbook/database"
	"net/http"
)

const MaximumDepth = 50

func sanitizeUserInput(r *http.Request) (*db.Task, []error) {

	err := r.ParseForm()
	if err != nil {
		return nil, []error{c.ErrCreateTaskParseForm}
	}

	var (
		content    = r.Form.Get("content")
		parentId   = r.Form.Get("parent_id")
		documentId = r.Form.Get("document_id")
	)

	if content == "" {
		return nil, []error{c.ErrCreateTaskEmptyContent}
	}

	if documentId == "" {
		return nil, []error{c.ErrCreateTaskEmptyDocumentId}
	} else if !v.IsValidUUID(documentId) {
		return nil, []error{c.ErrCreateTaskInvalidDocumentId}
	}

	if parentId == "" {
		return nil, []error{c.ErrCreateTaskEmptyParentId}
	} else if !v.IsValidUUID(parentId) {
		return nil, []error{c.ErrCreateTaskInvalidParentId}
	}

	task := db.Task{
		DocumentId: documentId,
		ParentId:   parentId,
		Content:    content,
	}
	return &task, nil
}

func createExecutor(r *http.Request) ([]db.Task, *e.Error) {
	var (
		task          *db.Task
		updated_tasks []db.Task
		errs          []error
	)

	task, errs = sanitizeUserInput(r)
	if errs != nil {
		return nil, e.New(`Check your inputs.`, errs)
	}

	// Call the db and make it official
	updated_tasks, errs = db.CreateTask(*task)
	if errs != nil {
		return nil, e.New(`Couldn't create task.`, errs)
	}

	return updated_tasks, nil
}

func Create(w http.ResponseWriter, r *http.Request) {
	task, errs := createExecutor(r)
	if errs != nil {
		c.ErrorHandler(w, r, errs)
	} else {
		c.SuccessHandler(w, r, task)
	}
}
