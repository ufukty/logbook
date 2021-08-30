package task

import (
	c "logbook/main/controller"
	db "logbook/main/database"
	"net/http"
	"strconv"
)

const MaximumDepth = 50

func sanitizeUserInput(w http.ResponseWriter, r *http.Request) (db.Task, error) {

	var (
		err error
	)

	err = r.ParseForm()
	if err != nil {
		c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrSterilizeUserInputParseForm, Underlying: err})
		return db.Task{}, err
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
		return db.Task{}, err
	}

	depth_int, err = strconv.Atoi(depth)
	if err != nil {
		c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrSterilizeUserInputDepthInt, Underlying: err})
		return db.Task{}, err
	}

	taskStatus_ts, err = db.StringToTaskStatus(taskStatus)
	if err != nil {
		c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrSterilizeUserInputTaskStatus, Underlying: err})
		return db.Task{}, err
	}

	task := db.Task{
		Content:     content,
		Degree:      degree_int,
		Depth:       depth_int,
		ParentId:    parentId,
		TaskGroupId: taskGroupId,
		TaskStatus:  taskStatus_ts,
	}
	return task, nil
}

func updateParentTask(w *http.ResponseWriter, r **http.Request, task *db.Task) error {

	var (
		parentTask db.Task
		err        error
	)

	if task.ParentId == "00000000-0000-0000-0000-000000000000" {
		return nil
	}

	parentTask, err = db.GetTaskByTaskId(task.ParentId)
	if err != nil {
		c.ErrorHandler(*w, *r, c.ControllerError{Wrapper: c.ErrParentCheck, Underlying: err})
		return err
	}

	// Increment the degree of parent,
	nextTasks, err := db.GetTaskByParentId(task.ParentId)
	if err != nil {
		c.ErrorHandler(*w, *r, c.ControllerError{Wrapper: c.ErrNextTaskCheck, Underlying: err})
		return err
	}
	totalDegree := 1
	for _, nextTask := range nextTasks {
		totalDegree += nextTask.Degree
	}
	parentTask.Degree = totalDegree

	// Call the db to update parent task to save changes
	_, err = db.UpdateTaskItem(parentTask)
	if err != nil {
		c.ErrorHandler(*w, *r, c.ControllerError{Wrapper: c.ErrCreateTaskUpdateParent, Underlying: err})
		return err
	}

	if parentTask.Depth >= MaximumDepth {
		return c.ErrMaximumDepthReached
	} else {
		return updateParentTask(w, r, &parentTask)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {

	var (
		task       db.Task
		parentTask db.Task
		err        error
	)

	task, err = sanitizeUserInput(w, r)
	if err != nil {
		return // sanitizeUserInput should already done logging and writing response
	}

	// Check if task group exists
	_, err = db.GetTaskGroupByTaskGroupId(task.TaskGroupId)
	if err != nil {
		c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrTaskGroupIdCheck, Underlying: err})
		return
	}

	// Call the db and make it official
	task, err = db.CreateTask(task)
	if err != nil {
		c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrCreateTaskCall, Underlying: err})
		return
	}

	if task.ParentId != "00000000-0000-0000-0000-000000000000" {
		parentTask, err = db.GetTaskByTaskId(task.ParentId)
		if err != nil {
			c.ErrorHandler(w, r, c.ControllerError{Wrapper: c.ErrParentCheck, Underlying: err})
			return
		}

		// Change status of parent to "drawer"
		if parentTask.TaskStatus == db.Active {
			parentTask.TaskStatus = db.Drawer
		}

		// Set the depth of child (current) task
		task.Depth = parentTask.Depth + 1

		// Check if the task is root or not,
		err = updateParentTask(&w, &r, &parentTask)
		if err != nil {
			return // updateParentTask should already done logging and writing response
		}
	}

	c.SuccessHandler(w, r, task)
}
