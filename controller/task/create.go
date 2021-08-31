package task

import (
	c "logbook/main/controller"
	db "logbook/main/database"
	"net/http"
	"strconv"
)

const MaximumDepth = 50

func sanitizeUserInput(r *http.Request) (db.Task, []error) {

	var (
		err  error
		errs []error
	)

	err = r.ParseForm()
	if err != nil {
		return db.Task{}, []error{err, c.ErrSterilizeUserInputParseForm}
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
		return db.Task{}, []error{err, c.ErrSterilizeUserInputDegreeInt}
	}

	depth_int, err = strconv.Atoi(depth)
	if err != nil {
		return db.Task{}, []error{err, c.ErrSterilizeUserInputDepthInt}
	}

	taskStatus_ts, errs = db.StringToTaskStatus(taskStatus)
	if errs != nil {
		return db.Task{}, append(errs, c.ErrSterilizeUserInputTaskStatus)
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

func updateParentTask(r **http.Request, task *db.Task) []error {

	var (
		parentTask db.Task
		errs       []error
	)

	if task.ParentId == "00000000-0000-0000-0000-000000000000" {
		return nil
	}

	parentTask, errs = db.GetTaskByTaskId(task.ParentId)
	if errs != nil {
		return append(errs, c.ErrUpdateParentParentCheck)
	}

	// Increment the degree of parent,
	nextTasks, errs := db.GetTaskByParentId(task.ParentId)
	if errs != nil {
		return append(errs, c.ErrUpdateParentNextTaskCheck)
	}
	totalDegree := 1
	for _, nextTask := range nextTasks {
		totalDegree += nextTask.Degree
	}
	parentTask.Degree = totalDegree

	// Call the db to update parent task to save changes
	_, errs = db.UpdateTaskItem(parentTask)
	if errs != nil {
		return append(errs, c.ErrUpdateParentSaveChanges)
	}

	if parentTask.Depth >= MaximumDepth {
		return []error{c.ErrUpdateParentMaximumDepthReached}
	} else {
		return updateParentTask(r, &parentTask)
	}
}

func createExecutor(r *http.Request) (db.Task, []error) {
	var (
		task       db.Task
		parentTask db.Task
		errs       []error
	)

	task, errs = sanitizeUserInput(r)
	if errs != nil {
		return db.Task{}, append(errs, c.ErrTaskCreateSanitize)
	}

	// Check if task group exists
	_, errs = db.GetTaskGroupByTaskGroupId(task.TaskGroupId)
	if errs != nil {
		return db.Task{}, append(errs, c.ErrTaskCreateTaskGroupIdCheck)
	}

	// Call the db and make it official
	task, errs = db.CreateTask(task)
	if errs != nil {
		return db.Task{}, append(errs, c.ErrTaskCreateCreateTaskCall)
	}

	if task.ParentId != "00000000-0000-0000-0000-000000000000" {
		parentTask, errs = db.GetTaskByTaskId(task.ParentId)
		if errs != nil {
			return db.Task{}, append(errs, c.ErrTaskCreateParentCheck)
		}

		// Change status of parent to "drawer"
		if parentTask.TaskStatus == db.Active {
			parentTask.TaskStatus = db.Drawer
		}

		// Set the depth of child (current) task
		task.Depth = parentTask.Depth + 1

		// Check if the task is root or not,
		errs = updateParentTask(&r, &parentTask)
		if errs != nil {
			return db.Task{}, append(errs, c.ErrTaskCreateUpdateParents)
		}
	}

	return task, nil
}

func Create(w http.ResponseWriter, r *http.Request) {
	task, errs := createExecutor(r)
	if errs != nil {
		c.ErrorHandler(w, r, errs)
	} else {
		c.SuccessHandler(w, r, task)
	}
}
