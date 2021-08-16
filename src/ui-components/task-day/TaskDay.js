import React from "react";

import "./TaskDay.css";

import TaskDayHeader from "./task-day-header/TaskDayHeader";
import TaskList from "./task-list/TaskList";

class TaskDay extends React.Component {
    constructor(props) {
        super();
        this.state = {
            day: props.data.day,
            tasks: props.data.tasks,
        };
    }

    render() {
        return (
            <div className="task-day">
                <TaskDayHeader className="task-day-header" />
                <TaskList tasks={this.state.tasks} />
            </div>
        );
    }
}

export default TaskDay;
