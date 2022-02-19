import React from "react";

import "./TaskList.css";

import Task from "./task/Task";

class TaskList extends React.Component {
    constructor(props) {
        super();
        this.state = {
            tasks: props.tasks,
            documentViewMode: props.documentViewMode,
        };
    }

    render() {
        var children = this.state.tasks.map((task) => (
            <Task key={task.task_id} task={task} documentViewMode={this.state.documentViewMode} />
        ));
        return <div className="task-list">{children}</div>;
    }
}

export default TaskList;
