import React from "react";

import "./TaskDay.css";

import TaskDayHeader from "./task-day-header/TaskDayHeader";
import TaskList from "./task-list/TaskList";

function TaskDay() {
    return (
        <div className="task-day">
            <TaskDayHeader className="task-day-header" />
            <TaskList />
        </div>
    );
}

export default TaskDay;
