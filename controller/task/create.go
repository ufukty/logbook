package task

import (
	c "logbook/main/controller"
	db "logbook/main/database"
	"net/http"
	"strconv"
)

func sanitizeUserInput(w http.ResponseWriter, r *http.Request) db.Task {

	var (
		err error
	)

	err = r.ParseForm()
	if err != nil {
		c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrSterilizeUserInputParseForm, Underlying: err})
	}

	var (
		content       = r.Form.Get("content")
		degree        = r.Form.Get("degree")
		depth         = r.Form.Get("depth")
		parentId      = r.Form.Get("parent_id")
		taskGroupId   = r.Form.Get("task_group_id")
		taskStatus    = r.Form.Get("task_status")
		degree_int    int
		depth_int     int
		taskStatus_ts db.TaskStatus
	)

	degree_int, err = strconv.Atoi(degree)
	if err != nil {
		c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrSterilizeUserInputDegreeInt, Underlying: err})
	}

	depth_int, err = strconv.Atoi(depth)
	if err != nil {
		c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrSterilizeUserInputDepthInt, Underlying: err})
	}

	taskStatus_ts, err = db.StringToTaskStatus(taskStatus)
	if err != nil {
		c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrSterilizeUserInputTaskStatus, Underlying: err})
	}

	task := db.Task{
		Content:     content,
		Degree:      degree_int,
		Depth:       depth_int,
		ParentId:    parentId,
		TaskGroupId: taskGroupId,
		TaskStatus:  taskStatus_ts,
	}
	return task
}

func Create(w http.ResponseWriter, r *http.Request) {

	task := sanitizeUserInput(w, r)

	if task.ParentId != "00000000-0000-0000-0000-000000000000" {
		_, err := db.GetTaskByTaskId(task.ParentId)
		if err != nil {
			c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrTaskGroupIdCheck, Underlying: err})
		}
	}

	_, err := db.GetTaskGroupByTaskGroupId(task.TaskGroupId)
	if err != nil {
		c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrTaskGroupIdCheck, Underlying: err})
	}

	task, err = db.CreateTask(task)
	if err != nil {
		c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrCreateTaskCall, Underlying: err})
	}

	c.SuccessHandler(w, r, task)
}
